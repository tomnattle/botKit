package commonHTTP

type Request struct {
	H Header `json:"header"`
	R struct {
		W string      `json:"w"`
		C string      `json:"c"`
		M string      `json:"m"`
		P interface{} `json:"p"`
	} `json:"request"`
}

func MakeReq(usefulRequestPointer interface{}) *Request {
	ins := &Request{}
	ins.H = *commonHeader
	ins.R.P = usefulRequestPointer
	return ins
}

func MakeReqWithLogID(logID string, usefulRequestPointer interface{}) *Request {
	ins := &Request{}
	ins.H = *commonHeader
	ins.H.LogID = logID
	ins.R.P = usefulRequestPointer
	return ins
}
