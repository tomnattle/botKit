package botEcho

import (
	"fmt"
	"github.com/ifchange/botKit/logger"
	"net/http"
)

type Logger struct {
	format func(info string) string
}

func newLogger(logID string, req *http.Request) *Logger {
	return &Logger{
		format: func(info string) string {
			return fmt.Sprintf("LogID:%s Path:%s RemoteIP:%s Info:%v",
				logID, req.RequestURI, req.RemoteAddr, info)
		},
	}
}

func (log *Logger) Debugf(format string, v ...interface{}) {
	logger.Debugf(log.format(fmt.Sprintf(format, v...)))
}

func (log *Logger) Infof(format string, v ...interface{}) {
	logger.Infof(log.format(fmt.Sprintf(format, v...)))
}

func (log *Logger) Warnf(format string, v ...interface{}) {
	logger.Warnf(log.format(fmt.Sprintf(format, v...)))
}

func (log *Logger) Errorf(format string, v ...interface{}) {
	logger.Errorf(log.format(fmt.Sprintf(format, v...)))
}
