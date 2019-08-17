package handler

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/SuperDogHuman/teraconnectgo/usecase"
	"github.com/labstack/echo/v4"
)

func getAuthoringLesson(c echo.Context) error {
	id := c.Param("id")

	ids := []string{id}
	if !IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lesson, err := usecase.GetAuthoringLesson(c.Request(), id)
	if err != nil {
		lessonErr, ok := err.(usecase.AuthoringLessonErrorCode)
		if ok && lessonErr == usecase.AuthoringLessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lesson)
}

func postAuthoringLesson(c echo.Context) error {
	lesson := new(domain.Lesson)

	if err := c.Bind(lesson); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := usecase.CreateAuthoringLesson(c.Request(), lesson); err != nil {
		fatalLog(err)
		authoringLessonErr, ok := err.(usecase.AuthoringLessonErrorCode)
		if ok && authoringLessonErr == usecase.InvalidAuthoringLessonParams {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, lesson)
}

func patchAuthoringLesson(c echo.Context) error {
	id := c.Param("id")

	ids := []string{id}
	// TODO add checking of avatarID, graphicIDs
	if !IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		fatalLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lesson, err := usecase.UpdateAuthoringLesson(id, c.Request())
	if err != nil {
		fatalLog(err)
		authoringLessonErr, ok := err.(usecase.AuthoringLessonErrorCode)
		if ok && authoringLessonErr == usecase.AuthoringLessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok && authoringLessonErr == usecase.InvalidAuthoringLessonParams {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lesson)
}
