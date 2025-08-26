package utils

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

func CustomErrHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var code int
	var msg string

	if _, isUrlError := err.(*url.Error); isUrlError {
		code = http.StatusBadGateway
		msg = "Bad Gateway: " + err.Error()
	} else {
		code = http.StatusInternalServerError
		msg = err.Error()
	}

	c.String(code, msg)
}