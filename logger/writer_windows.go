// +build windows

package logger

import (
	"io"
	"os"
)

func GetOutput() io.Writer { return os.Stdout }
