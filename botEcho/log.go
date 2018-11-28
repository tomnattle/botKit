package botEcho

import (
	"net/http"
)

type Logger struct{}

func newLogger(logID string, req *http.Request) *Logger {
	return &Logger{}
}
