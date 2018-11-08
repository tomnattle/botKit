package admin

import (
	"bytes"
	"github.com/ifchange/botKit/config"
	"github.com/ifchange/botKit/signature"
	"io"
	"net/http"
)

var (
	cfg *config.URIConfig
)

func init() {
	cfg = config.GetConfig().URI
	if cfg == nil {
		panic("botKit-admin products config is nil")
	}
}

func getURI() string {
	return cfg.Admin
}

func AdminPOST(subURI string) (*http.Request, io.Writer, error) {
	body := &bytes.Buffer{}
	req, err := http.NewRequest("POST", getURI()+subURI, body)
	if err != nil {
		return nil, nil, err
	}
	err = signature.AddSignature(req)
	if err != nil {
		return nil, nil, err
	}
	return req, body, nil
}
