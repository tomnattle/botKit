package botEcho

import (
	"github.com/labstack/echo"
)

type HandlerFunc func(Context) error

func (ins *Server) CONNECT(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return ins.Echo.CONNECT(path, handler(h), m...)
}

func (ins *Server) DELETE(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return ins.Echo.DELETE(path, handler(h), m...)
}

func (ins *Server) GET(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return ins.Echo.GET(path, handler(h), m...)
}

func (ins *Server) HEAD(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return ins.Echo.HEAD(path, handler(h), m...)
}

func (ins *Server) OPTIONS(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return ins.Echo.OPTIONS(path, handler(h), m...)
}

func (ins *Server) PATCH(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return ins.Echo.PATCH(path, handler(h), m...)
}

func (ins *Server) POST(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return ins.Echo.POST(path, handler(h), m...)
}

func (ins *Server) PUT(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return ins.Echo.PUT(path, handler(h), m...)
}

func (ins *Server) TRACE(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return ins.Echo.TRACE(path, handler(h), m...)
}

func (ins *Server) Any(path string, h HandlerFunc, m ...echo.MiddlewareFunc) []*echo.Route {
	return ins.Echo.Any(path, handler(h), m...)
}
