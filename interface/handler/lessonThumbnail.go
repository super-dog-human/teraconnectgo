package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func postLessonThumbnail(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	url, err := usecase.CreateLessonThumbnailBlankFile(c.Request(), id)
	if err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := response{url}

	return c.JSON(http.StatusCreated, response)
}

type response struct {
	URL string `json:"url"`
}
