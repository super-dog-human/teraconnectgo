package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func postIntroductionLesson(c echo.Context) error {
	if err := usecase.CreateIntroductionLesson(c.Request()); err != nil {
		fatalLog(err)
		LessonErr, ok := err.(usecase.LessonErrorCode)
		if ok && LessonErr == usecase.InvalidLessonParams {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "succeeded")
}
