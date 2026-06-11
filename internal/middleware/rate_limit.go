package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/swaindhruti/pharmastock-backend/internal/common"
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
			return common.APIErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		},
		DenyHandler: func(c *echo.Context, identifier string, err error) error {
			return common.APIErrorResponse(c, http.StatusTooManyRequests, "Too many requests. Please try again later.")
		},
	}

	return middleware.RateLimiterWithConfig(config)
}
