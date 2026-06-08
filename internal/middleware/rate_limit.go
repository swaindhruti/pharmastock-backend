package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func RateLimit(requests int, duration time.Duration) echo.MiddlewareFunc {

	if requests <= 0 {
		requests = 60
	}

	if duration <= 0 {
		duration = time.Minute
	}

	rate := float64(requests) / float64(duration.Seconds())

	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate, Burst: requests, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(c *echo.Context) (string, error) {
			id := c.RealIP()
			return id, nil
		},
		ErrorHandler: func(c *echo.Context, err error) error {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"success": false,
				"error":   "Internal Server Error",
				"data":    nil,
			})
		},
		DenyHandler: func(c *echo.Context, identifier string, err error) error {
			return c.JSON(http.StatusTooManyRequests, map[string]any{
				"success": false,
				"error":   "Too many requests. Please try again later.",
				"data":    nil,
			})
		},
	}

	return middleware.RateLimiterWithConfig(config)
}
