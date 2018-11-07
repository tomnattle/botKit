package signature

import (
	"bytes"
	"net/http"
	"testing"
)

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
	signatureStr, err := Signature(timeStamp, nonce, cfg.SecretKey)
	if err != nil {
		t.Fatal(err)
	}
	if signatureStr != signature {
		t.Fatal("not equal")
	}
}
