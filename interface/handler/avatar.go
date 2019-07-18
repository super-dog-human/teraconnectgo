package handler

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/SuperDogHuman/teraconnectgo/usecase"
	"github.com/labstack/echo/v4"
)

func getAvatars(c echo.Context) error {
	avatars, err := usecase.GetAvailableAvatars(c.Request())

	if err != nil {
		authErr, ok := err.(domain.AuthErrorCode)
		if ok && authErr == domain.UserNotFound {
			// when token is valid but user account not exists.
			return c.JSON(http.StatusNotFound, err.Error())
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

func createAvatars(c echo.Context) error {
	objectRequest := new(domain.StorageObjectRequest)
	if err := c.Bind(objectRequest); err != nil {
		fatalLog(err)
		authErr, ok := err.(domain.AuthErrorCode)
		if ok && authErr == domain.UserNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	signedURLs, err := usecase.CreateAvatarsAndBlankFile(c.Request(), *objectRequest)
	if err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, signedURLs)
}
