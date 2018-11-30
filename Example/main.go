package main

import (
	"github.com/ifchange/botKit/botEcho"
	"github.com/ifchange/botKit/commonHTTP"
)

func main() {
	e := botEcho.New()
	e.POST("/", index)
	e.Run()
}

func index(c botEcho.Context) error {
	rsp := commonHTTP.MakeRsp(nil)
	c.Logger().Debugf("Fi")
	return c.JSON(rsp)
}
