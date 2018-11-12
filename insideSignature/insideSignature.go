package insideSignature

import (
	"fmt"
	"github.com/ifchange/botKit/config"
	"net/http"
)

var cfgInsideSignature string

func init() {
	cfg := config.GetConfig()
	if cfg == nil {
		panic("signature config error")
	}
	cfgInsideSignature = cfg.InsideSignature
}

func AddSignature(req *http.Request) error {
	if req == nil {
		return fmt.Errorf("botKit signature nil http-request")
	}
	req.Header.Add("x", cfgInsideSignature)
	return nil
}

func VerifySignature(req *http.Request) (pass bool, err error) {
	signatureStr := req.Header.Get("x")
	if signatureStr == "" {
		err = fmt.Errorf("signature can not be null")
		return
	}
	if cfgInsideSignature == signatureStr {
		pass = true
		return
	}
	return
}
