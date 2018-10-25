// +build windows

package grace

import (
	"net/http"
)

func Serve(s *http.Server) error { return s.ListenAndServe() }
