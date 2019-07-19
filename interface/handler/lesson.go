package handler

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/SuperDogHuman/teraconnectgo/usecase"
	"github.com/labstack/echo/v4"
)

func getLessons(c echo.Context) error {
	// TODO add pagination
	return c.JSON(http.StatusOK, "")
}

func getLesson(c echo.Context) error {
	id := c.Param("id")

	ids := []string{id}
	if !IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lesson, err := usecase.GetAvailableLesson(c.Request(), id)
	if err != nil {
		lessonErr, ok := err.(usecase.LessonErrorCode)
		if ok && lessonErr == usecase.LessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok && lessonErr == usecase.LessonNotAvailable {
			warnLog(lessonErr)
			return c.JSON(http.StatusForbidden, err.Error())
		}
		authErr, ok := err.(domain.AuthErrorCode)
		if ok {
			warnLog(authErr)
			return c.JSON(http.StatusUnauthorized, err.Error())
		}

		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lesson)
}

func deleteLesson(c echo.Context) error {
	id := c.Param("id")

	ids := []string{id}
	if !IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		fatalLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	if err := usecase.DeleteOwnLessonById(c.Request(), id); err != nil {
		fatalLog(err)
		lessonErr, ok := err.(usecase.LessonErrorCode)
		if ok && lessonErr == usecase.LessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok && lessonErr == usecase.LessonNotAvailable {
			return c.JSON(http.StatusForbidden, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "the lesson has deleted.")
}
