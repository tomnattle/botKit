package botEcho

import (
	"github.com/labstack/echo"
)

type HandlerFunc func(Context) error

func (ins *Server) POST(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return ins.Echo.POST(path, handler(h), m...)
}

func handler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return h(Context{c})
	}
}
