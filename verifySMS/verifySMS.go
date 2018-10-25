package verifySMS

import (
	"fmt"
	"ifchange/tsketch/kit/Redis"
	"ifchange/tsketch/kit/SMS"
	"ifchange/tsketch/kit/config"
	"math/rand"
	"time"
)

var (
	contentModel string
	// default 1 minute
	time2live time.Duration = time.Duration(1) * time.Minute
)

func init() {
	contentModel = config.GetConfig().SMS.Model
	time2live = time.Minute *
		time.Duration(config.GetConfig().SMS.DurationMinutes)
}

func VerifyAuthCode(phone, authCode string) (pass bool, err error) {
	redis, err := Redis.GetRedis()
	if err != nil {
		return false, err
	}
	defer redis.Close()

	saveAuthCode, ok := getAuthCode(redis, phone)
	if !ok {
		return false, nil
	}
	if saveAuthCode != authCode {
		return false, nil
	}
	return true, nil
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
	env := *config.GetEnvironment()
	if env == config.TEST || env == config.DEV {
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

func checkSendDuration(redis *Redis.RedisCommon, phone string) error {
	result, err := redis.Cmd("EXISTS", key(phone)).Int()
	if err != nil {
		return fmt.Errorf("kit-verifySMS checkSendDuration error %v", err)
	}
	if result > 0 {
		return fmt.Errorf("kit-verifySMS authCode is exist")
	}
	return nil
}

func saveAuthCode(redis *Redis.RedisCommon, phone, authCode string) error {
	err := redis.Cmd("SETEX", key(phone), int(time2live.Seconds()), authCode).Err
	if err != nil {
		return fmt.Errorf("exec redis query error %v", err)
	}
	return nil
}

func getAuthCode(redis *Redis.RedisCommon, phone string) (string, bool) {
	str, err := redis.Cmd("GET", key(phone)).Str()
	if err != nil {
		return str, false
	}
	return str, true
}

func key(phone string) string {
	return Redis.FormatKey(fmt.Sprintf("verifySMS_Phone_%s", phone))
}

func makeContent(authCode string) string {
	return fmt.Sprintf(contentModel, authCode, int(time2live.Minutes()))
}

func authCodeGenerator() string {
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", newRand.Int31n(1000000))
}
