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
	"net/url"
	"strconv"
	"time"
)

const (
	ConstTimeFormat = "20060102150405"
)

type Session struct {
	SrcID     int    `json:"src_id"`
	ManagerID int    `json:"manager_id"`
	Expire    string `json:"expire"`
	Args      string `json:"args"`
	Signature string `json:"signature"`
}

func GenerateSession(srcID int, managerID int, duration time.Duration, argsMapping map[string]string) (string, error) {
	if argsMapping == nil {
		argsMapping = make(map[string]string)
	}
	args := marshalArgs(argsMapping)

	expire := time.Now().Add(duration)
	secretKey, err := getSecretKey(managerID)
	if err != nil {
		return "", fmt.Errorf("srcID:%d managerID:%d getSecretKey error %v",
			srcID, managerID, err)
	}
	return NewSession(srcID, managerID, expire, args, secretKey)
}

func VerifySession(session string) (srcID int, managerID int, args map[string]string, pass bool, err error) {
	jsonSource, err := base64.URLEncoding.DecodeString(session)
	if err != nil {
		return srcID, managerID, args, pass, err
	}
	s := &Session{}
	err = json.Unmarshal(jsonSource, s)
	if err != nil {
		return srcID, managerID, args, pass, err
	}
	srcID = s.SrcID
	managerID = s.ManagerID

	expireTime, err := time.Parse(ConstTimeFormat, s.Expire)
	if err != nil {
		return srcID, managerID, args, pass, fmt.Errorf("VerifySession parse expire %v error %v", s.Expire, err)
	}

	if expireTime.Before(time.Now()) {
		pass = false
		return
	}

	secretKey, err := getSecretKey(s.ManagerID)
	if err != nil {
		return srcID, managerID, args, pass, fmt.Errorf("VerifySession srcID:%d managerID:%d getSecretKey error %v",
			s.SrcID, s.ManagerID, err)
	}

	newSession, err := NewSession(srcID, managerID, expireTime, s.Args, secretKey)
	if err != nil {
		return srcID, managerID, args, pass, fmt.Errorf("VerifySession srcID:%d managerID:%d NewSession error %v",
			s.SrcID, s.ManagerID, err)
	}

	if newSession != session {
		pass = false
		return
	}

	args, err = unmarshalArgs(s.Args)
	if err != nil {
		return srcID, managerID, args, pass, err
	}
	pass = true
	return
}

func NewSession(srcID int, managerID int, expire time.Time, args string, secretKey string) (string, error) {
	expireStr := expire.Format(ConstTimeFormat)
	source := strconv.Itoa(srcID) + strconv.Itoa(managerID) + expireStr + args + secretKey
	sha1er := sha1.New()
	io.WriteString(sha1er, source)
	signature := fmt.Sprintf("%x", sha1er.Sum(nil))

	data, err := json.Marshal(&Session{
		SrcID:     srcID,
		ManagerID: managerID,
		Expire:    expireStr,
		Args:      args,
		Signature: signature,
	})
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

func getSecretKey(managerID int) (string, error) {
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

func marshalArgs(args map[string]string) string {
	urlValues := url.Values{}
	for k, v := range args {
		urlValues.Set(k, v)
	}
	return urlValues.Encode()
}

func unmarshalArgs(args string) (map[string]string, error) {
	if len(args) == 0 {
		return nil, nil
	}
	urlValues, err := url.ParseQuery(args)
	if err != nil {
		return nil, fmt.Errorf("try parse url encoding error %v %v", args, err)
	}
	result := make(map[string]string)
	for k := range urlValues {
		result[k] = urlValues.Get(k)
	}
	return result, nil
}
