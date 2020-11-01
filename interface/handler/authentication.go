package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/domain"
)

// Authentication validates JWT in header.
func Authentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			if _, err := domain.ValidTokenClaims(request); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token.")
			}
			return next(c)
		}
	}
}
