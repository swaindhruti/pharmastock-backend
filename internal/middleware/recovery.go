package middleware

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func Recovery() echo.MiddlewareFunc {
	return middleware.Recover()
}
