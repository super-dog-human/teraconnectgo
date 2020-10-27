package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func getLessons(c echo.Context) error {
	lessons, err := usecase.GetLessonsByConditions(c.Request())
	if err != nil {
		lessonErr, ok := err.(usecase.LessonErrorCode)
		if ok && lessonErr == usecase.LessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok {
			fatalLog(lessonErr)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		authErr, ok := err.(domain.AuthErrorCode)
		if ok {
			fatalLog(authErr)
			return c.JSON(http.StatusUnauthorized, err.Error())
		}

		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lessons)
}

func getLesson(c echo.Context) error {
	id := c.Param("id")

	ids := []string{id}
	if !IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	var lesson domain.Lesson
	var err error
	if c.Param("for_authoring") == "true" {
		lesson, err = usecase.GetPrivateLesson(c.Request(), id)
	} else {
		lesson, err = usecase.GetPublicLesson(c.Request(), id)
	}

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

func postLesson(c echo.Context) error {
	lesson := new(domain.Lesson)

	if err := c.Bind(lesson); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := usecase.CreateLesson(c.Request(), lesson); err != nil {
		fatalLog(err)
		LessonErr, ok := err.(usecase.LessonErrorCode)
		if ok && LessonErr == usecase.InvalidLessonParams {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, lesson)
}

func patchLesson(c echo.Context) error {
	id := c.Param("id")

	ids := []string{id}
	// TODO add checking of avatarID, graphicIDs
	if !IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		fatalLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lesson, err := usecase.UpdateLesson(id, c.Request())
	if err != nil {
		fatalLog(err)
		LessonErr, ok := err.(usecase.LessonErrorCode)
		if ok && LessonErr == usecase.LessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok && LessonErr == usecase.InvalidLessonParams {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
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

	if err := usecase.DeleteOwnLessonByID(c.Request(), id); err != nil {
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
