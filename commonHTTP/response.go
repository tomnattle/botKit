package commonHTTP

import (
	"fmt"
	"strings"
)

type Response struct {
	H *Header `json:"header"`
	R struct {
		ErrNo   int         `json:"err_no"`
		ErrMsg  string      `json:"err_msg"`
		Results interface{} `json:"results"`
	} `json:"response"`
}

func MakeRsp(usefulResponsePointer interface{}) *Response {
	ins := &Response{}
	ins.H = commonHeader
	ins.R.Results = usefulResponsePointer
	return ins
}

func (ins *Response) Errorf(err error, errCode int, msg ...string) error {
	if err == nil {
		return nil
	}

	errMsg, ok := fullCodeMapping[errCode]
	if !ok {
		errMsg, ok = shortCodeMapping[errCode]
		if errMsg, ok = shortCodeMapping[errCode]; ok {
			errCode = shortCodeToFullCode(errCode)
		} else {
			errMsg = fmt.Sprintf("undefind err msg code %d", errCode)
		}
	}

	logMsg := errMsg + " - " + err.Error()

	if len(msg) > 0 {
		logMsg += " - "
		logMsg += strings.Join(msg, " ")
	}

	return &errCommon{
		errCode: errCode,
		errMsg:  errMsg,
		logMsg:  logMsg,
	}
}
