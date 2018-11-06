package commonHTTP

type Response struct {
	ErrNo   int         `json:"err_no"`
	ErrMsg  string      `json:"err_msg"`
	Results interface{} `json:"results"`
}
