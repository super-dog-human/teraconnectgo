package handler

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func getAvatars(c echo.Context) error {
	// TODO pagination.
	avatars, err := domain.GetAvailableAvatars(c.Request())

	if err != nil {
		fatalLog(errors.WithStack(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(avatars) == 0 {
		errMessage := "avatars not found"
		warnLog(errMessage)
		return c.JSON(http.StatusNotFound, errMessage)
	}

	return c.JSON(http.StatusOK, avatars)
}
