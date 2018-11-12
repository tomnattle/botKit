package middleware

import (
	"fmt"
	"github.com/ifchange/botKit/commonHTTP"
	"github.com/ifchange/botKit/insideSignature"
	"github.com/labstack/echo"
)

func Signature() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rsp := commonHTTP.MakeRsp(nil)
			pass, err := insideSignature.VerifySignature(c.Request())
			if err != nil {
				return rsp.Errorf(fmt.Errorf("verify signature err %v", err), 4001)
			}
			if pass {
				return next(c)
			}
			return rsp.Errorf(fmt.Errorf("signature wrong"), 4005)
		}
	}
}
