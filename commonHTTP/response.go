package commonHTTP

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type Response struct {
	H *Header `json:"header"`
	R struct {
		ErrNo   int         `json:"err_no"`
		ErrMsg  string      `json:"err_msg"`
		Results interface{} `json:"results"`
	} `json:"response"`
}

func GetRsp(reader io.Reader, usefulResponsePointer interface{}) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("read from io.Reader error %v", err)
	}
	ins := &Response{}
	ins.R.Results = usefulResponsePointer
	err = json.Unmarshal(data, ins)
	if err != nil {
		return fmt.Errorf("json unmarshal error %v -- %v -- %v",
			err, usefulResponsePointer, string(data))
	}
	if ins.R.ErrNo == 0 {
		return nil
	}
	return fmt.Errorf("errNo:%v errMsg:%v", ins.R.ErrNo, ins.R.ErrMsg)
}

func MakeRsp(usefulResponsePointer interface{}) *Response {
	ins := &Response{}
	ins.H = commonHeader
	ins.R.Results = usefulResponsePointer
	return ins
}
