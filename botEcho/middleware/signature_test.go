package middleware

import (
	"bytes"
	"github.com/ifchange/botKit/util"
	"net/http"
	"testing"
	"time"
)

func TestSignature(t *testing.T) {
	//timeStamp : 2018101817
	timeStamp := time.Now().Format("2006010215")
	nonce := util.RandStr(15)
	signature, err := creatSignature(timeStamp, nonce, cfg.SecretKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("timeStamp =%v nonce = %v signature = %v", timeStamp, nonce, signature)
}

func TestAddSignature(t *testing.T) {
	req, err := http.NewRequest("POST", "127.0.0.1", &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	err = AddSignature(req)
	if err != nil {
		t.Fatal(err)
	}
	timeStamp := req.Header.Get("timeStamp")
	signature := req.Header.Get("signature")
	nonce := req.Header.Get("nonce")
	signatureStr, err := creatSignature(timeStamp, nonce, cfg.SecretKey)
	if err != nil {
		t.Fatal(err)
	}
	if signatureStr != signature {
		t.Fatal("not equal")
	}
}
