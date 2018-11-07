package session

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"github.com/ifchange/botKit/signature"
	"github.com/ifchange/botKit/util"
	"time"
)

const (
	salt string = "lle3sx2wECsSEvs2a23vx92cvS3z12E9c1F"
)

type Session struct {
	Args      string `json:"args"`
	Nonce     string `json:"nonce"`
	TimeStamp string `json:"timeStamp"`
	Signature string `json:"signature"`
}

func GenerateSession(argsStructPointer interface{}) (session string, err error) {
	writer := new(bytes.Buffer)
	err = gob.NewEncoder(writer).Encode(argsStructPointer)
	if err != nil {
		return session, err
	}
	args := writer.String()
	timeStamp := time.Now().Format("2006010215")
	nonce := util.RandStr(15)
	signature, err := signature.Signature(timeStamp, nonce, args)
	if err != nil {
		return session, err
	}
	data, err := json.Marshal(&Session{
		Args:      args,
		Nonce:     nonce,
		TimeStamp: timeStamp,
		Signature: signature,
	})
	if err != nil {
		return session, err
	}
	session = base64.URLEncoding.EncodeToString(data)
	return
}

func VerifySession(session string, argsStructPointer interface{}) (pass bool, err error) {
	jsonSource, err := base64.URLEncoding.DecodeString(session)
	if err != nil {
		return pass, err
	}
	s := &Session{}
	err = json.Unmarshal(jsonSource, s)
	if err != nil {
		return pass, err
	}
	signature, err := signature.Signature(s.TimeStamp, s.Nonce, s.Args)
	if err != nil {
		return pass, err
	}
	if signature != s.Signature {
		pass = false
		return
	}
	err = gob.NewDecoder(bytes.NewBufferString(s.Args)).Decode(argsStructPointer)
	if err != nil {
		return pass, err
	}
	pass = true
	return
}
