package middleware

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func RequestID() echo.MiddlewareFunc {
	return middleware.RequestID()
}
