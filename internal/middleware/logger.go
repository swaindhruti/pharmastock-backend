package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type statusCodeCapture struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusCodeCapture) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func Logger(logger *zap.Logger) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			wrapped := &statusCodeCapture{ResponseWriter: c.Response()}
			c.SetResponse(wrapped)

			err := next(c)

			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = "unknown"
			}

			fields := []zap.Field{
				zap.String("request_id", requestID),
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.Int("status", wrapped.statusCode),
				zap.Duration("latency", time.Since(start)),
				zap.String("client_ip", c.RealIP()),
				zap.String("user_agent", c.Request().UserAgent()),
			}

			if err != nil {
				fields = append(fields, zap.Error(err))
				logger.Error("http_request", fields...)
			} else {
				logger.Info("http_request", fields...)
			}

			return err
		}
	}
}
