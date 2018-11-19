// inside
package signature

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
	req.Header.Add("X", cfgInsideSignature)
	return nil
}

func VerifySignature(req *http.Request) error {
	signatureStr := req.Header.Get("X")
	if signatureStr == "" {
		return fmt.Errorf("VerifySignature signature can not be null")
	}
	if cfgInsideSignature == signatureStr {
		return nil
	}
	return fmt.Errorf("VerifySignature unauthorized")
}
