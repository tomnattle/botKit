package verifySMS

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ifchange/botKit/Redis"
	"github.com/ifchange/botKit/SMS"
	"github.com/ifchange/botKit/config"
	"github.com/ifchange/botKit/logger"
	"math/rand"
	"time"
)

var (
	contentModel string
	// default 1 minute
	sendDuration time.Duration = time.Duration(1) * time.Minute
	// default 3 minute
	authCodeLiveDuration time.Duration = time.Duration(3) * time.Minute
)

var (
	ErrAuthCodeIsExist = errors.New("authCode is exist")
)

func init() {
	contentModel = config.GetConfig().SMS.Model
	sendDuration = time.Minute *
		time.Duration(config.GetConfig().SMS.DurationMinutes)
	authCodeLiveDuration = time.Minute *
		time.Duration(config.GetConfig().SMS.AuthCodeDurationMinutes)
}

func SendVerifySMS(phone string) error {
	redis, err := Redis.GetRedis()
	if err != nil {
		return err
	}
	defer redis.Close()

	if err := checkSendDuration(redis, phone); err != nil {
		return err
	}
	authCode := authCodeGenerator()
	env := config.GetConfig().Environment
	if env == "test" || env == "dev" {
		phoneRunes := []rune(phone)
		authCodeRunes := make([]rune, 6)
		authCodeRunes[0] = phoneRunes[0]
		authCodeRunes[1] = phoneRunes[2]
		authCodeRunes[2] = phoneRunes[4]
		authCodeRunes[3] = phoneRunes[6]
		authCodeRunes[4] = phoneRunes[8]
		authCodeRunes[5] = phoneRunes[10]
		authCode = string(authCodeRunes)
	}
	content := makeContent(authCode)
	if err := SMS.Request(phone, content); err != nil {
		return err
	}
	err = saveAuthCode(redis, phone, authCode)
	if err != nil {
		return err
	}
	return nil
}

func VerifyAuthCode(phone, authCode string) (pass bool, err error) {
	redis, err := Redis.GetRedis()
	if err != nil {
		return false, err
	}
	defer redis.Close()

	allAuthCode, ok := getAuthCode(redis, phone)
	if !ok {
		return false, nil
	}
	now := time.Now()
	for _, saveAuthCode := range allAuthCode {
		if saveAuthCode.Expire.Sub(now) < 0 {
			continue
		}
		if saveAuthCode.AuthCode != authCode {
			continue
		}
		err := deleteAuthCode(redis, phone)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func checkSendDuration(redis *Redis.RedisCommon, phone string) error {
	result, err := redis.Cmd("EXISTS", sendDurationKey(phone)).Int()
	if err != nil {
		return fmt.Errorf("kit-verifySMS checkSendDuration error %v", err)
	}
	if result > 0 {
		return ErrAuthCodeIsExist
	}
	return nil
}

type AuthCodeSaver struct {
	AuthCode string    `json:"auth_code"`
	Expire   time.Time `json:"expire"`
}

func getAuthCode(redis *Redis.RedisCommon, phone string) ([]*AuthCodeSaver, bool) {
	all, err := redis.Cmd("GET", authCodeSaverKey(phone)).Bytes()
	if err != nil {
		return nil, false
	}
	allAuthCode := []*AuthCodeSaver{}
	if err := json.Unmarshal(all, &allAuthCode); err != nil {
		logger.Warnf("getAuthCode Unmarshal json error %v", err)
		return nil, false
	}
	return allAuthCode, true
}

func saveAuthCode(redis *Redis.RedisCommon, phone, authCode string) error {
	newAuthCode := &AuthCodeSaver{
		AuthCode: authCode,
		Expire:   time.Now().Add(authCodeLiveDuration),
	}
	all, ok := getAuthCode(redis, phone)
	if ok {
		all = append(all, newAuthCode)
	} else {
		all = []*AuthCodeSaver{newAuthCode}
	}
	data, err := json.Marshal(&all)
	if err != nil {
		return fmt.Errorf("saveAuthCode Marshal json error %v", err)
	}
	err = redis.Cmd("SETEX", authCodeSaverKey(phone), int(authCodeLiveDuration.Seconds()), data).Err
	if err != nil {
		return fmt.Errorf("exec redis query error %v", err)
	}
	err = redis.Cmd("SETEX", sendDurationKey(phone), int(sendDuration.Seconds()), authCode).Err
	if err != nil {
		return fmt.Errorf("exec redis query error %v", err)
	}
	return nil
}

func deleteAuthCode(redis *Redis.RedisCommon, phone string) error {
	if err := redis.Cmd("DEL", authCodeSaverKey(phone)).Err; err != nil {
		return fmt.Errorf("exec redis query error %v", err)
	}
	if err := redis.Cmd("DEL", sendDurationKey(phone)).Err; err != nil {
		return fmt.Errorf("exec redis query error %v", err)
	}
	return nil
}

func sendDurationKey(phone string) string {
	return Redis.FormatKey(fmt.Sprintf("verifySMS_sendDuration_Phone_%s", phone))
}

func authCodeSaverKey(phone string) string {
	return Redis.FormatKey(fmt.Sprintf("verifySMS_authCodeSaver_Phone_%s", phone))
}

func makeContent(authCode string) string {
	return fmt.Sprintf(contentModel, authCode, int(sendDuration.Minutes()))
}

func authCodeGenerator() string {
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", newRand.Int31n(1000000))
}
