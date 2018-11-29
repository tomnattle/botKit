package botEcho

import (
	"encoding/json"
	"fmt"
	"github.com/ifchange/botKit/commonHTTP"
	"github.com/ifchange/botKit/config"
	"github.com/ifchange/botKit/util"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"runtime"
)

type Context struct {
	// inside echo context
	echo.Context
	// link to http request and response
	request  *http.Request
	response **interface{}
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
		c.response = new(*interface{})
		// unmarshal
		reply := commonHTTP.MakeRsp(nil)
		requestBody, err := ioutil.ReadAll(c.request.Body)
		if err != nil {
			return reply.Errorf(
				fmt.Errorf("common handler read request body error %v", err),
				4001)
		}
		commonRequest := &struct {
			H commonHTTP.Header `json:"header"`
			R struct {
				W string           `json:"w"`
				C string           `json:"c"`
				M string           `json:"m"`
				P *json.RawMessage `json:"p"`
			} `json:"request"`
		}{}
		err = json.Unmarshal(requestBody, commonRequest)
		if err != nil {
			return reply.Errorf(
				fmt.Errorf("common handler unmarshal request body error %v", err),
				4001)
		}
		c.CommonHeader = &commonRequest.H
		c.W = commonRequest.R.W
		c.C = commonRequest.R.C
		c.M = commonRequest.R.M
		c.P = commonRequest.R.P

		if len(c.CommonHeader.LogID) < 5 {
			c.CommonHeader.LogID = util.RandStr(22)
		}
		c.logger = newLogger(c.CommonHeader.LogID, c.request)

		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				stack := make([]byte, 4<<12)
				length := runtime.Stack(stack, true)
				errHandler(fmt.Errorf("[PANIC RECOVER] %v %s", err, stack[:length]), c)
			}
		}()

		c.Logger().Infof("request info %s", string(requestBody))
		if err := h(c); err != nil {
			errHandler(err, c)
			return nil
		}
		if *c.response == nil {
			c.Logger().Infof("nil response")
			return nil
		}
		response, err := json.Marshal(**c.response)
		if err != nil {
			c.Logger().Infof("response info unknown response %v", c.response)
			return nil
		}
		c.Logger().Infof("response info %s", string(response))
		return nil
	}
}

func (c Context) Logger() *Logger {
	return c.logger
}

func (c Context) GetReq(usefulRequestPointer interface{}) error {
	if c.P == nil {
		return fmt.Errorf("botEcho.Context.GetReq error empty body")
	}
	err := json.Unmarshal([]byte(*c.P), usefulRequestPointer)
	if err != nil {
		return fmt.Errorf("botEcho.Context.GetReq error json unmarshal error %v %v", err, usefulRequestPointer)
	}
	return nil
}

func (c Context) JSON(response interface{}) error {
	*c.response = &response
	return c.Context.JSON(http.StatusOK, response)
}

func ErrHandler(err error, echoC echo.Context) {
	// make context
	c := Context{Context: echoC}
	c.request = echoC.Request()
	c.logger = newLogger("unknown", c.request)
	errHandler(err, c)
}

func errHandler(err error, c Context) {
	var (
		code   = http.StatusOK
		rsp    = commonHTTP.MakeRsp(nil)
		logMsg = ""
	)

	if errC, ok := err.(*commonHTTP.ErrCommon); ok {
		switch config.GetConfig().Environment {
		case "dev":
			rsp.R.ErrNo = errC.ErrCode
			rsp.R.ErrMsg = errC.ErrMsg + errC.LogMsg
			logMsg = errC.LogMsg
		default:
			rsp.R.ErrNo = errC.ErrCode
			rsp.R.ErrMsg = errC.ErrMsg
			logMsg = errC.LogMsg
		}
	} else {
		rsp.R.ErrNo = -1
		rsp.R.ErrMsg = "SYSTEM ERROR, please call backend ASAP"
		logMsg = err.Error()
	}

	c.Logger().Warnf("err:%v info:%v", rsp.R.ErrNo, logMsg)

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD { // echo Issue #608
			err = c.NoContent(code)
		} else {
			err = c.Context.JSON(code, rsp)
		}
		if err != nil {
			c.Logger().Errorf("errHandler %v", err)
		}
	}
}
