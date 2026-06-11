package common

import "github.com/labstack/echo/v5"

type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func APISuccessResponse(c *echo.Context, statusCode int, data any) error {
	return c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
	})
}

func APISuccessMessage(c *echo.Context, statusCode int, message string, data any) error {
	return c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func APIErrorResponse(c *echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}
