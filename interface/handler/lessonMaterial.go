package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func getLessonMaterials(c echo.Context) error {
	lessonID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lessonMaterial, err := usecase.GetLessonMaterial(c.Request(), lessonID)
	if err != nil {
		fatalLog(err)
		lessonErr, ok := err.(usecase.LessonErrorCode)
		if ok && lessonErr == usecase.LessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok && lessonErr == usecase.LessonNotAvailable {
			return c.JSON(http.StatusForbidden, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lessonMaterial)
}

func postLessonMaterial(c echo.Context) error {
	lessonID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lessonMaterial := new(domain.LessonMaterial)
	if err := c.Bind(lessonMaterial); err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := usecase.CreateLessonMaterial(c.Request(), lessonID, *lessonMaterial); err != nil {
		fatalLog(err)
		lessonErr, ok := err.(usecase.LessonErrorCode)
		if ok && lessonErr == usecase.LessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok && lessonErr == usecase.LessonNotAvailable {
			return c.JSON(http.StatusForbidden, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "succeeded")
}

func putLessonMaterial(c echo.Context) error {
	lessonID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lessonMaterial := new(domain.LessonMaterial)
	if err := c.Bind(lessonMaterial); err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := usecase.UpdateLessonMaterial(c.Request(), lessonID, *lessonMaterial); err != nil {
		fatalLog(err)
		lessonErr, ok := err.(usecase.LessonErrorCode)
		if ok && lessonErr == usecase.LessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok && lessonErr == usecase.LessonNotAvailable {
			return c.JSON(http.StatusForbidden, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "succeeded")
}
