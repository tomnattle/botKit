package admin

import (
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

func GetAdminURI() string {
	return cfg.Admin
}
