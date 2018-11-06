// +build !windows

package grace

import (
	"github.com/facebookgo/grace/gracehttp"
	"net/http"
)

func Serve(s *http.Server) error { return gracehttp.Serve(s) }
