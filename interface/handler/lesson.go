package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	var lesson domain.Lesson
	var err error

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	if c.QueryParam("for_authoring") == "true" {
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

func getCurrentUserLessons(c echo.Context) error {
	lessons, err := usecase.GetCurrentUserLessons(c.Request())

	if err != nil {
		fatalLog(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if len(lessons) == 0 {
		return c.JSON(http.StatusNotFound, "lesson doesn't exist.")
	}

	return c.JSON(http.StatusOK, lessons)
}

func postLesson(c echo.Context) error {
	params := new(usecase.NewLessonParams)
	lesson := new(domain.Lesson)

	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// TODO validate newLesson

	if err := usecase.CreateLesson(c.Request(), params, lesson); err != nil {
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errMessage := "Invalid lessonID error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	var params map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	needsCopyThumbnail := c.QueryParam("move_thumbnail") == "true"
	if err := usecase.UpdateLessonWithMaterial(id, c.Request(), needsCopyThumbnail, &params); err != nil {
		fatalLog(err)
		LessonErr, ok := err.(usecase.LessonErrorCode)
		if ok && LessonErr == usecase.LessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok && LessonErr == usecase.InvalidLessonParams {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "succeeded.")
}
