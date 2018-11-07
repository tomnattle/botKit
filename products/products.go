package products

import (
	"fmt"
	"github.com/ifchange/botKit/config"
)

var (
	cfg *config.URIConfig
)

func init() {
	cfg = config.GetConfig().URI
	if cfg == nil {
		panic("botKit products config is nil")
	}
}

func GetURI(productID int) (uri string, err error) {
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
