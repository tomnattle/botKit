package botEcho

import (
	"github.com/labstack/echo"
)

type Context struct {
	echo.Context
}

func handler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return h(Context{c})
	}
}
