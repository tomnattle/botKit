package session

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const (
	ConstFromA = "A"
	ConstFromB = "B"
	ConstFromC = "C"

	constMaxDuration = time.Duration(7*24) * time.Hour
)

type Session struct {
	From      string `json:"from"` // A B C
	SrcID     int    `json:"src_id"`
	ManagerID int    `json:"manager_id"`
	UserID    int    `json:"user_id"`
	Expire    int    `json:"expire"`
	Signature string `json:"signature"`
}

func GenerateSession(from string, srcID, managerID, userID int, duration time.Duration,
	getSecretKey func(srcID int, managerID int) (secretKey string, err error)) (string, error) {
	if duration > constMaxDuration {
		return "", fmt.Errorf("GenerateSession error out of max expire:%v want:%v",
			constMaxDuration, duration)
	}
	expire := time.Now().Unix() + int64(duration.Seconds())
	secretKey, err := getSecretKey(srcID, managerID)
	if err != nil {
		return "", fmt.Errorf("GenerateSession from:%s srcID:%d managerID:%d userID:%d getSecretKey error %v",
			from, srcID, managerID, userID, err)
	}
	return newSession(from, srcID, managerID, userID, expire, secretKey)
}

func VerifySession(from string, session string,
	getSecretKey func(srcID int, managerID int) (secretKey string, err error)) (*Session, error) {
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
	now := time.Now().Unix()
	if now > int64(s.Expire) {
		return nil, fmt.Errorf("VerifySession session is timeout")
	}
	if int64(s.Expire)-now > int64(constMaxDuration.Seconds()) {
		return nil, fmt.Errorf("VerifySession error out of max expire %v",
			s.Expire)
	}
	secretKey, err := getSecretKey(s.SrcID, s.ManagerID)
	if err != nil {
		return nil, fmt.Errorf("VerifySession srcID:%d managerID:%d userID:%d getSecretKey error %v",
			s.SrcID, s.ManagerID, s.UserID, err)
	}
	newSignature, err := newSignature(s.From, s.SrcID, s.ManagerID, s.UserID, int64(s.Expire), secretKey)
	if err != nil {
		return nil, fmt.Errorf("VerifySession srcID:%d managerID:%d userID:%d NewSignature error %v",
			s.SrcID, s.ManagerID, s.UserID, err)
	}
	if newSignature != s.Signature {
		return nil, fmt.Errorf("VerifySession unauthorized")
	}
	return s, nil
}

func newSession(from string, srcID, managerID, userID int, expireSeconds int64, secretKey string) (string, error) {
	signature, err := newSignature(from, srcID, managerID, userID, expireSeconds, secretKey)
	if err != nil {
		return "", fmt.Errorf("newSignature error %v", err)
	}
	data, err := json.Marshal(&Session{
		From:      from,
		SrcID:     srcID,
		ManagerID: managerID,
		UserID:    userID,
		Expire:    int(expireSeconds),
		Signature: signature,
	})
	if err != nil {
		return "", fmt.Errorf("newSession json marshal error %v", err)
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

func newSignature(from string, srcID, managerID, userID int, expireSeconds int64, secretKey string) (string, error) {
	switch from {
	case ConstFromA, ConstFromB, ConstFromC:
	default:
		return "", fmt.Errorf("newSignature unknown from %s", from)
	}

	source := fmt.Sprintf("%s%d%d%d%d%s", from, srcID, managerID, userID, expireSeconds, secretKey)
	sha1er := sha1.New()
	io.WriteString(sha1er, source)
	return fmt.Sprintf("%x", sha1er.Sum(nil)), nil
}
