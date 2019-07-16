package handler

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/labstack/echo/v4"
)

func getAvatars(c echo.Context) error {
	avatars, err := domain.GetAvailableAvatars(c.Request())

	if err != nil {
		authErr, ok := err.(domain.AuthErrorCode)
		if ok && authErr == domain.UserNotFound {
			// when token is valid but user account not exists.
			return c.JSON(http.StatusNotFound, authErr)
		} else {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	if len(avatars) == 0 {
		errMessage := "avatars not found"
		warnLog(errMessage)
		return c.JSON(http.StatusNotFound, errMessage)
	}

	return c.JSON(http.StatusOK, avatars)
}
