package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func postUserThumbnail(c echo.Context) error {
	url, err := usecase.CreateUserThumbnailBlankFile(c.Request())
	if err != nil {
		if ok := errors.Is(err, usecase.UserNotAvailable); ok {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := thumbnailResponse{url}
	return c.JSON(http.StatusCreated, response)
}
