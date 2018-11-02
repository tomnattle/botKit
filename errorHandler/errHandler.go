package errorHandler

import (
	"fmt"
	"github.com/ifchange/botKit/config"
	"github.com/ifchange/botKit/logger"
	"github.com/labstack/echo"
	"net/http"
	"runtime/debug"
	"strings"
)

var codeMapping map[int]string

func init() {
	codeMapping = make(map[int]string)
	for _, errConfig := range errCodeConfig() {
		code := config.GetConfig().AppID*1000000 + config.GetConfig().SubAppID*10000 + errConfig.code
		codeMapping[code] = errConfig.msg
	}
	logger.Printf("all errCode config %v\n", codeMapping)
}

type errConfig struct {
	code int
	msg  string
}

type ErrCode struct {
	Code int    `json:"err_no"`
	Msg  string `json:"err_msg"`
}

func (err *ErrCode) String() string {
	return fmt.Sprintf("errNo:%v errMsg:%v",
		err.Code, err.Msg)
}

type errCommon struct {
	errCode int
	errMsg  string
	logMsg  string
}

func (err *errCommon) Error() string {
	if err == nil {
		return "nil errCommon in errHandler package"
	}
	return err.errMsg
}

func (ins *ErrCode) Errorf(err error, errCode int, msg ...string) error {
	if err == nil {
		return nil
	}

	errMsg, ok := codeMapping[errCode]
	if !ok {
		errMsg = fmt.Sprintf("undefind err msg code %d", errCode)
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

func ErrHandler(err error, c echo.Context) {
	var (
		code   = http.StatusOK
		msg    = &ErrCode{}
		logMsg = ""
	)

	if errC, ok := err.(*errCommon); ok {
		switch config.GetEnvironment() {
		case config.DEV, config.TEST:
			msg.Code = errC.errCode
			msg.Msg = errC.errMsg + errC.logMsg
			logMsg = errC.logMsg
		default:
			msg.Code = errC.errCode
			msg.Msg = errC.errMsg
			logMsg = errC.logMsg
		}
	} else {
		msg.Code = -1
		msg.Msg = "SYSTEM ERROR, please call backend ASAP"
		logMsg = err.Error()
	}

	c.Logger().Warnf("uri:%s err:%v info:%v stack:%s", c.Request().RequestURI, msg.Code, logMsg,
		func() string { return string(debug.Stack()) }())

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD { // echo Issue #608
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
