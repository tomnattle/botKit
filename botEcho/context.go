package botEcho

import (
	"github.com/labstack/echo"
)

type Context struct {
	echo.Context
	commonRequest  bool
	commonResponse bool
}
