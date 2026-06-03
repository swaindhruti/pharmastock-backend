package common

import "github.com/labstack/echo/v5"

func APISuccessResponse(
	c *echo.Context,
	statusCode int,
	Data any,
) error

func APIErrorResponse(
	c *echo.Context,
	statusCode int,
	message string,
) error
