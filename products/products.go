package products

import (
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
	case 3:
		uri = cfg.WinMode
	case 4:
		uri = cfg.Tsketch
	default:
		err = fmt.Errorf("botKit products error productID:%d is not defind",
			productID)
	}
	return
}

func getName(productID int) (name string, err error) {
	switch productID {
	case 1:
		name = "面试bot"
	case 3:
		name = "决胜力"
	case 4:
		name = "人才画像"
	case 6:
		name = "與情BI"
	case 7:
		name = "情商"
	case 8:
		name = "岗位评估"
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

func GetProduct(productID int) (*Product, error) {
	name, err := getName(productID)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:   productID,
		Name: name,
	}, nil
}

func ProductPOST(productID int, subURI string, body io.Reader) (*http.Request, error) {
	basicURI, err := getURI(productID)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", basicURI+subURI, body)
	if err != nil {
		return nil, err
	}
	err = signature.AddSignature(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
