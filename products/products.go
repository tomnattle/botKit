package products

import (
	"bytes"
	"fmt"
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
		panic("botKit-products products config is nil")
	}
}

func getURI(productID int) (uri string, err error) {
	switch productID {
	case 1:
		uri = cfg.ChatBot
	case 2:
		uri = cfg.WinMode
	case 3:
		uri = cfg.Tsketch
	default:
		err = fmt.Errorf("botKit products error productID:%d is not defind",
			productID)
	}
	return
}

type Product struct {
	ID   int    `json:"product_id"`
	Name string `json:"product_name"`
}

func ProductPOST(productID int, subURI string) (*http.Request, io.Writer, error) {
	body := &bytes.Buffer{}
	basicURI, err := getURI(productID)
	if err != nil {
		return nil, nil, err
	}
	req, err := http.NewRequest("POST", basicURI+subURI, body)
	if err != nil {
		return nil, nil, err
	}
	err = signature.AddSignature(req)
	if err != nil {
		return nil, nil, err
	}
	return req, body, nil
}
