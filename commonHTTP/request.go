package commonHTTP

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type Request struct {
	H *Header `json:"header"`
	R struct {
		W string      `json:"w"`
		C string      `json:"c"`
		M string      `json:"m"`
		P interface{} `json:"p"`
	} `json:"request"`
}

func GetReq(reader io.Reader, usefulRequestPointer interface{}) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("read from io.Reader error %v", err)
	}
	ins := &Request{}
	ins.R.P = usefulRequestPointer
	err = json.Unmarshal(data, ins)
	if err != nil {
		return fmt.Errorf("json unmarshal error %v %v", err, usefulRequestPointer)
	}
	return nil
}

func MakeReq(usefulRequestPointer interface{}) *Request {
	ins := &Request{}
	ins.H = commonHeader
	ins.R.P = usefulRequestPointer
	return ins
}
