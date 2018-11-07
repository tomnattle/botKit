package middleware

import (
	"github.com/ifchange/botKit/util"
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
