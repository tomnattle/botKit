package botEcho

import (
	"fmt"
	"github.com/ifchange/botKit/botEcho/grace"
	"github.com/ifchange/botKit/commonHTTP"
	"github.com/ifchange/botKit/config"
	"github.com/ifchange/botKit/logger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

type Server struct {
	*echo.Echo
}

func New() *Server {
	appName := config.GetConfig().AppName
	// load config
	fmt.Printf("%s start, run environment is %s\n",
		appName, config.GetConfig().Environment)
	cfg := config.GetConfig()
	if cfg == nil {
		panic(fmt.Sprintf("%s start fail, load config error", appName))
	}
	fmt.Printf("%s start, listen port is %s\n",
		appName, cfg.Addr)
	// init service
	e := &Server{Echo: echo.New()}
	// init log
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logger.GetOutput(),
	}))
	if config.GetConfig().Environment == "dev" {
		e.Logger.SetLevel(log.DEBUG)
	} else {
		e.Logger.SetLevel(log.WARN)
	}
	e.Logger.SetOutput(logger.GetOutput())

	// middleware
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("2M"))
	// common error handler
	e.HTTPErrorHandler = commonHTTP.ErrHandler

	e.Server.Addr = cfg.Addr
	return e
}

func (ins *Server) Run() {
	err := grace.Serve(ins.Echo.Server)
	if err != nil {
		panic(err)
	}
}
