package commonHTTP

type Request struct {
	W string      `json:"w"`
	C string      `json:"c"`
	M string      `json:"m"`
	P interface{} `json:"p"`
}
