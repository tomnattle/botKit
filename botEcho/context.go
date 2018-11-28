package botEcho

import (
	"encoding/json"
	"fmt"
	"github.com/ifchange/botKit/commonHTTP"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
)

type Context struct {
	// inside echo context
	echo.Context
	// link to http request
	request *http.Request
	// log in each request
	logger *Logger
	// common request
	CommonHeader *commonHTTP.Header
	W            string
	C            string
	M            string
	P            *json.RawMessage
}

func handler(h HandlerFunc) echo.HandlerFunc {
	return func(echoC echo.Context) error {
		c := Context{Context: echoC}
		c.request = echoC.Request()

		reply := commonHTTP.MakeRsp(nil)
		body, err := ioutil.ReadAll(c.request.Body)
		if err != nil {
			return reply.Errorf(
				fmt.Errorf("common handler read request body error %v", err),
				4001)
		}
		commonRequest := &struct {
			H *commonHTTP.Header `json:"header"`
			R struct {
				W string           `json:"w"`
				C string           `json:"c"`
				M string           `json:"m"`
				P *json.RawMessage `json:"p"`
			} `json:"request"`
		}{}
		err = json.Unmarshal(body, commonRequest)
		if err != nil {
			return reply.Errorf(
				fmt.Errorf("common handler unmarshal request body error %v", err),
				4001)
		}
		c.CommonHeader = commonRequest.H
		c.W = commonRequest.R.W
		c.C = commonRequest.R.C
		c.M = commonRequest.R.M
		c.P = commonRequest.R.P

		c.logger = newLogger(c.CommonHeader.LogID, c.request)
		return h(c)
	}
}

func (c Context) Logger() *Logger {
	return c.logger
}

func (c Context) Request() *http.Request {
	return c.request
}
