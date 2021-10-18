package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// Authentication validates JWT in header.
func Authentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			if _, err := domain.ValidTokenClaims(request); err != nil {
				log.Printf("failed to token validation.\n")
				log.Printf("%v\n", errors.WithStack(err).Error())
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token.")
			}
			return next(c)
		}
	}
}

func CSRFTokenCookie() echo.MiddlewareFunc {
	return echo.WrapMiddleware(csrf.Protect(
		csrfInitialString(),
		csrf.Path("/"),
		csrf.TrustedOrigins([]string{infrastructure.OriginURL()}),
	))
}

func CSRFTokenHeader() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Response().Header()
			header.Set(echo.HeaderAccessControlExposeHeaders, "x-csrf-token")
			header.Set(echo.HeaderXCSRFToken, csrf.Token(c.Request()))
			return next(c)
		}
	}
}

func csrfInitialString() []byte {
	if infrastructure.AppEnv() == "development" {
		return []byte("32-byte-long-auth-key")
	} else {
		return []byte(os.Getenv("COOKIE_SECRET"))
	}
}
