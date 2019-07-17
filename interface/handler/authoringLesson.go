package handler

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
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

	lesson, err := domain.GetAuthoringLesson(c.Request(), id)
	if err != nil {
		lessonErr, ok := err.(domain.AuthoringLessonErrorCode)
		if ok && lessonErr == domain.AuthoringLessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lesson)
}

func createAuthoringLesson(c echo.Context) error {
	postedLesson := new(domain.Lesson)

	if err := c.Bind(postedLesson); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	lesson, err := domain.CreateAuthoringLesson(c.Request(), *postedLesson)
	if err != nil {
		fatalLog(err)

		authoringLessonErr, ok := err.(domain.AuthoringLessonErrorCode)
		if ok && authoringLessonErr == domain.InvalidAuthoringLessonParams {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, lesson)
}

func updateAuthoringLesson(c echo.Context) error {
	id := c.Param("id")

	ids := []string{id}
	// TODO add checking of avatarID, graphicIDs
	if !IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		fatalLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lesson, err := domain.UpdateAuthoringLesson(id, c.Request())
	if err != nil {
		fatalLog(err)

		authoringLessonErr, ok := err.(domain.AuthoringLessonErrorCode)
		if ok && authoringLessonErr == domain.AuthoringLessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok && authoringLessonErr == domain.InvalidAuthoringLessonParams {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lesson)
}
