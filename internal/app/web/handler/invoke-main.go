package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func InvokeMain(c echo.Context) error {
	return c.Render(http.StatusOK, "field-enter.html", echo.Map{})
}
