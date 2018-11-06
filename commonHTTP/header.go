package commonHTTP

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
