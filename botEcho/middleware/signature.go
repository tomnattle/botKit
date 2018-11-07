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
			timeStamp := c.Request().Header.Get("timeStamp")
			signatureStr := c.Request().Header.Get("signature")
			nonce := c.Request().Header.Get("nonce")
			if signatureStr == "" {
				return rsp.Errorf(fmt.Errorf("signature can not be null"), 4005)
			}
			selfSignatureStr, err := signature.Signature(timeStamp, nonce, cfg.SecretKey)
			if err != nil {
				return rsp.Errorf(fmt.Errorf("signature create err %v", err), 4005)
			}
			if selfSignatureStr == signatureStr {
				return next(c)
			}
			return rsp.Errorf(fmt.Errorf("signature wrong"), 4005)
		}
	}
}
