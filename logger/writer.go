// +build !windows

package logger

import (
	"io"
)

func GetOutput() io.Writer { return writer }
