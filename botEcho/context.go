package botEcho

import (
	"encoding/json"
	"fmt"
	"github.com/ifchange/botKit/commonHTTP"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"runtime"
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
		// make context
		c := Context{Context: echoC}
		c.request = echoC.Request()
		// unmarshal
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

		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				stack := make([]byte, 4<<10)
				length := runtime.Stack(stack, true)
				c.Logger().Printf("[PANIC RECOVER] %v %s", err, stack[:length])
				c.Context.Error(err)
			}
		}()

		if err := h(c); err == nil {
			return err
		}
		c.Logger.Printf("Logger error err")
		return h(c)
	}
}

func (c Context) Logger() *Logger {
	return c.logger
}

func (c Context) Request() *http.Request {
	return c.request
}
