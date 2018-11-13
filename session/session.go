// outside
package session

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ifchange/botKit/admin"
	"github.com/ifchange/botKit/commonHTTP"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	ConstTimeFormat = "20060102150405"
	ConstFromA      = "A"
	ConstFromB      = "B"
	ConstFromC      = "C"
)

type Session struct {
	From      string `json:"From"` // A B C
	SrcID     int    `json:"src_id"`
	ManagerID int    `json:"manager_id"`
	UserID    int    `json:"user_id"`
	Expire    string `json:"expire"`
	Signature string `json:"signature"`
}

func GenerateSession(from string, srcID, managerID, userID int, duration time.Duration) (string, error) {
	expire := time.Now().Add(duration)
	secretKey, err := GetSecretKey(managerID)
	if err != nil {
		return "", fmt.Errorf("GenerateSession from:%s srcID:%d managerID:%d userID:%d getSecretKey error %v",
			from, srcID, managerID, userID, err)
	}
	return NewSession(from, srcID, managerID, userID, expire, secretKey)
}

func VerifySession(from string, session string) (*Session, error) {
	jsonSource, err := base64.URLEncoding.DecodeString(session)
	if err != nil {
		return nil, fmt.Errorf("VerifySession base64 decode error %v", err)
	}
	s := &Session{}
	err = json.Unmarshal(jsonSource, s)
	if err != nil {
		return nil, fmt.Errorf("VerifySession json unmarshal error %v", err)
	}
	if from != s.From {
		return nil, fmt.Errorf("VerifySession diff from %s:%s", from, s.From)
	}
	expireTime, err := time.Parse(ConstTimeFormat, s.Expire)
	if err != nil {
		return nil, fmt.Errorf("VerifySession parse expire %v error %v", s.Expire, err)
	}
	if expireTime.Before(time.Now()) {
		return nil, fmt.Errorf("VerifySession session is timeout")
	}
	secretKey, err := GetSecretKey(s.ManagerID)
	if err != nil {
		return nil, fmt.Errorf("VerifySession srcID:%d managerID:%d userID:%d getSecretKey error %v",
			s.SrcID, s.ManagerID, s.UserID, err)
	}
	newSession, err := NewSession(s.From, s.SrcID, s.ManagerID, s.UserID, expireTime, secretKey)
	if err != nil {
		return nil, fmt.Errorf("VerifySession srcID:%d managerID:%d userID:%d NewSession error %v",
			s.SrcID, s.ManagerID, s.UserID, err)
	}
	if newSession != session {
		return nil, fmt.Errorf("VerifySession unauthorized")
	}
	return s, nil
}

func NewSession(from string, srcID, managerID, userID int, expire time.Time, secretKey string) (string, error) {
	switch from {
	case ConstFromA, ConstFromB, ConstFromC:
	default:
		return "", fmt.Errorf("NewSession unknown from %s", from)
	}

	expireStr := expire.Format(ConstTimeFormat)
	source := strconv.Itoa(srcID) + strconv.Itoa(managerID) + strconv.Itoa(userID) + expireStr + secretKey
	sha1er := sha1.New()
	io.WriteString(sha1er, source)
	signature := fmt.Sprintf("%x", sha1er.Sum(nil))

	data, err := json.Marshal(&Session{
		From:      from,
		SrcID:     srcID,
		ManagerID: managerID,
		UserID:    userID,
		Expire:    expireStr,
		Signature: signature,
	})
	if err != nil {
		return "", fmt.Errorf("NewSession json marshal error %v", err)
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

func GetSecretKey(managerID int) (string, error) {
	body := &bytes.Buffer{}
	reqBody := commonHTTP.MakeReq(&struct {
		ManagerID int `json:"id"`
	}{managerID})

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	_, err = body.Write(reqData)
	if err != nil {
		return "", fmt.Errorf("try write body error %v", err)
	}

	req, err := admin.AdminPOST("/companies/getsecretkey", body)
	if err != nil {
		return "", fmt.Errorf("admin make request error %v", err)
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("admin request error %v", err)
	}
	defer rsp.Body.Close()

	secretKey := ""

	err = commonHTTP.GetRsp(rsp.Body, &secretKey)
	if err != nil {
		return "", err
	}

	if len(secretKey) == 0 {
		return "", fmt.Errorf("empty secretKey")
	}
	return secretKey, nil
}
