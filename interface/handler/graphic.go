package handler

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/usecase"
	"github.com/labstack/echo/v4"
)

func getGraphics(c echo.Context) error {
	// TODO pagination.
	graphics, err := usecase.GetAvailableGraphics(c.Request())

	if err != nil {
		fatalLog(err)
		authErr, ok := err.(domain.AuthErrorCode)
		if ok && authErr == domain.UserNotFound {
			// when token is valid but user account not exists.
			return c.JSON(http.StatusNotFound, err.Error())
		} else {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	if len(graphics) == 0 {
		errMessage := "graphics not found"
		warnLog(errMessage)
		return c.JSON(http.StatusNotFound, errMessage)
	}

	return c.JSON(http.StatusOK, graphics)
}

func postGraphics(c echo.Context) error {
	objectRequest := new(domain.StorageObjectRequest)
	if err := c.Bind(objectRequest); err != nil {
		fatalLog(err)
		authErr, ok := err.(domain.AuthErrorCode)
		if ok && authErr == domain.UserNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	signedURLs, err := usecase.CreateGraphicsAndBlankFile(c.Request(), *objectRequest)
	if err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, signedURLs)
}
