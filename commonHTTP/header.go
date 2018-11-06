package commonHTTP

import (
	"github.com/ifchange/botKit/config"
)

var commonHeader *Header

func init() {
	cfg := config.GetConfig()
	if cfg == nil {
		panic("commonHTTP nil config")
	}
	commonHeader = &Header{
		AppID:    cfg.AppID,
		LogID:    cfg.Environment,
		Provider: cfg.AppName,
	}
}

type Header struct {
	AppID    int    `json:"appid"`
	LogID    string `json:"log_id"`
	UID      string `json:"uid"`
	UName    string `json:"uname"`
	Provider string `json:"provider"`
	SignID   string `json:"signid"`
	Version  string `json:"version"`
	IP       string `json:"ip"`
}
