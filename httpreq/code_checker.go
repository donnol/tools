package httpreq

import (
	"fmt"
	"net/http"
)

type (
	CodeChecker func(code int) error
)

func CodeIs200(code int) error {
	if code != http.StatusOK {
		return fmt.Errorf("bad http code: %d", code)
	}
	return nil
}
