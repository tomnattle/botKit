package middleware

import (
	"fmt"
	"github.com/ifchange/botKit/commonHTTP"
	"github.com/ifchange/botKit/config"
	"github.com/ifchange/botKit/signature"
	"github.com/labstack/echo"
)

var cfg *config.SignatureConfig

func init() {
	cfg = config.GetConfig().Signature
	if cfg == nil {
		panic("signature config error")
	}
}

func Signature() echo.MiddlewareFunc {
	return signatureWithConfig(*cfg)
}

func signatureWithConfig(config config.SignatureConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rsp := commonHTTP.MakeRsp(nil)
			pass, err := signature.VerifySignature(c.Request())
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
