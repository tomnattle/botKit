package commonHTTP

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type Request struct {
	H Header `json:"header"`
	R struct {
		W string      `json:"w"`
		C string      `json:"c"`
		M string      `json:"m"`
		P interface{} `json:"p"`
	} `json:"request"`
}

func GetReq(reader io.Reader, emptyUsefulRequestPointer interface{}) (usefulRequestPointer interface{}, err error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read from io.Reader error %v", err)
	}
	ins := &Request{}
	ins.R.P = emptyUsefulRequestPointer
	err = json.Unmarshal(data, ins)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal error %v %v", err, emptyUsefulRequestPointer)
	}
	return ins.R.P, nil
}
