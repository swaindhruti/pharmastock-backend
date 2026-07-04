package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/swaindhruti/pharmastock-backend/internal/common"
)

type contextKey string

const (
	ContextUserID      contextKey = "user_id"
	ContextUserRole    contextKey = "user_role"
	ContextReferenceID contextKey = "reference_id"
)

func AuthRequired(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return common.APIErrorResponse(c, http.StatusUnauthorized, "authorization header required")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				return common.APIErrorResponse(c, http.StatusUnauthorized, "invalid authorization header format")
			}

			claims, err := validateJWT(parts[1], jwtSecret)
			if err != nil {
				return common.APIErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
			}

			c.Set(string(ContextUserID), claims.UserID)
			c.Set(string(ContextUserRole), claims.Role)
			c.Set(string(ContextReferenceID), claims.ReferenceID)

			return next(c)
		}
	}
}

func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			role, ok := c.Get(string(ContextUserRole)).(string)
			if !ok {
				return common.APIErrorResponse(c, http.StatusForbidden, "access denied")
			}

			for _, allowed := range roles {
				if role == allowed {
					return next(c)
				}
			}

			return common.APIErrorResponse(c, http.StatusForbidden, "insufficient permissions")
		}
	}
}

func GetUserID(c *echo.Context) int64 {
	id, _ := c.Get(string(ContextUserID)).(int64)
	return id
}

func GetUserRole(c *echo.Context) string {
	role, _ := c.Get(string(ContextUserRole)).(string)
	return role
}

func GetReferenceID(c *echo.Context) int64 {
	ref, _ := c.Get(string(ContextReferenceID)).(int64)
	return ref
}
