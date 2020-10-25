package handler

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/usecase"
	"github.com/labstack/echo/v4"
)

func putLessonPack(c echo.Context) error {
	id := c.Param("id")

	ids := []string{id}
	if !IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		fatalLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	if err := usecase.PackLesson(c.Request(), id); err != nil {
		fatalLog(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, "succeeded")
}
