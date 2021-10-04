package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

type getLessonShortResponse struct {
	NextCursor string               `json:"nextCursor"`
	Lessons    []domain.ShortLesson `json:"lessons"`
}

func getLessons(c echo.Context) error {
	categoryID, err := strconv.ParseInt(c.QueryParam("category_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if categoryID == 0 {
		return c.JSON(http.StatusNotFound, "category_id is blank.")
	}

	cursorStr := c.QueryParam("next_cursor")

	lessons, nextCursorStr, err := usecase.GetLessonsByCategoryID(c.Request(), categoryID, cursorStr)
	if err != nil {
		lessonErr, ok := err.(usecase.LessonErrorCode)
		if ok && lessonErr == usecase.LessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok {
			fatalLog(lessonErr)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := getLessonShortResponse{Lessons: lessons, NextCursor: nextCursorStr}

	return c.JSON(http.StatusOK, response)
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
		viewKey := c.QueryParam("view_key")
		lesson, err = usecase.GetPublicLesson(c.Request(), id, viewKey)
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
		return c.JSON(http.StatusNotFound, "lesson doesn't exist")
	}

	return c.JSON(http.StatusOK, lessons)
}

func postLesson(c echo.Context) error {
	params := new(usecase.NewLessonParams)

	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if !params.IsIntroduction {
		if params.Title == "" {
			return c.JSON(http.StatusBadRequest, "title is blank")
		} else if params.SubjectID == 0 {
			return c.JSON(http.StatusBadRequest, "subjectID is blank")
		} else if params.JapaneseCategoryID == 0 {
			return c.JSON(http.StatusBadRequest, "japanseCategoryID is blank")
		}
	}

	lesson := new(domain.Lesson)
	var err error
	if params.IsIntroduction {
		err = usecase.CreateIntroductionLesson(c.Request(), lesson)
	} else {
		err = usecase.CreateLesson(c.Request(), params, lesson)
	}

	if err != nil {
		fatalLog(err)
		if ok := errors.Is(err, usecase.InvalidLessonParams); ok {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if ok := errors.Is(err, domain.AlreadyIntroductionExist); ok {
			return c.JSON(http.StatusConflict, err.Error())
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
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	if err := usecase.UpdateLessonWithMaterial(id, c.Request(), needsCopyThumbnail, requestID, &params); err != nil {
		fatalLog(err)
		LessonErr, ok := err.(usecase.LessonErrorCode)
		if ok && LessonErr == usecase.LessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok && LessonErr == usecase.InvalidLessonParams {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "succeeded")
}
