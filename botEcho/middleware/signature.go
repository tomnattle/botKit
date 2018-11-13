package middleware

import (
	"fmt"
	"github.com/ifchange/botKit/commonHTTP"
	"github.com/ifchange/botKit/signature"
	"github.com/labstack/echo"
)

func Signature() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rsp := commonHTTP.MakeRsp(nil)
			err := signature.VerifySignature(c.Request())
			if err != nil {
				return rsp.Errorf(fmt.Errorf("verify signature err %v", err), 4001)
			}
			return next(c)
		}
	}
}
