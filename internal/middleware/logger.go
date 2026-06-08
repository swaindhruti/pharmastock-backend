package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}

// statusCodeCapture wraps http.ResponseWriter to capture the status code
type statusCodeCapture struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusCodeCapture) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func Logger(service string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			wrapped := &statusCodeCapture{ResponseWriter: c.Response()}
			c.SetResponse(wrapped)

			RequestID := c.Response().Header().Get(echo.HeaderXRequestID)
			if RequestID == "" {
				RequestID = "unknown"
			}

			err := next(c)
			latency := time.Since(start)

			logger.Info(
				"http_request",
				zap.String("service", service),
				zap.String("request_id", RequestID),
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.Int("status", wrapped.statusCode),
				zap.Duration("latency", latency),
				zap.String("client_ip", c.RealIP()),
				zap.String("user_agent", c.Request().UserAgent()),
			)

			return err
		}
	}
}
