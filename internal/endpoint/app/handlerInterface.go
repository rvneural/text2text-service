package app

import "github.com/labstack/echo/v4"

type HandleFunc interface {
	HandleRequest(c echo.Context) error
}
