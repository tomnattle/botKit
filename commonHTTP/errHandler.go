package commonHTTP

import (
	"fmt"
	"github.com/ifchange/botKit/config"
	"strings"
)

var (
	fullCodeMapping  map[int]string
	shortCodeMapping map[int]string
	fullCodeBasic    int
)

func shortCodeToFullCode(shortCode int) (fullCode int) {
	fullCode = fullCodeBasic + shortCode
	return
}

func init() {
	fullCodeBasic = config.GetConfig().AppID*1000000 + config.GetConfig().SubAppID*10000
	fullCodeMapping = make(map[int]string)
	shortCodeMapping = make(map[int]string)
	for _, errConfig := range errCodeConfig() {
		fullCodeMapping[shortCodeToFullCode(errConfig.code)] = errConfig.msg
		shortCodeMapping[errConfig.code] = errConfig.msg
	}
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

	return &ErrCommon{
		ErrCode: errCode,
		ErrMsg:  errMsg,
		LogMsg:  logMsg,
	}
}

type errConfig struct {
	code int
	msg  string
}

type ErrCommon struct {
	ErrCode int
	ErrMsg  string
	LogMsg  string
}

func (err *ErrCommon) Error() string {
	if err == nil {
		return "nil errCommon in errHandler package"
	}
	return err.ErrMsg
}
