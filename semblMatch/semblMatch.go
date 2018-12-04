package semblMatch

import (
	"github.com/ifchange/botKit/config"
	"github.com/ifchange/botKit/signature"
	"io"
	"net/http"
)

var (
	cfg *config.NLPConfig
)

func init() {
	cfg = config.GetConfig().NLP
	if cfg == nil || cfg.SemblMatchServer == "" {
		panic("botKit-admin products config is nil")
	}
}

func getURI() string {
	return cfg.SemblMatchServer
}

func SemblMatchPOST(subURI string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest("POST", getURI()+subURI, body)
	if err != nil {
		return nil, err
	}
	err = signature.AddSignature(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
