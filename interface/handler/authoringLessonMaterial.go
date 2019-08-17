package handler

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/SuperDogHuman/teraconnectgo/usecase"
	"github.com/labstack/echo/v4"
)

func getAuthoringLessonMaterials(c echo.Context) error {
	lessonID := c.Param("id")

	ids := []string{lessonID}
	if !IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		fatalLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lessonMaterial, err := usecase.GetAuthoringLessonMaterial(c.Request(), lessonID)
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

func postAuthoringLessonMaterial(c echo.Context) error {
	lessonID := c.Param("id")

	ids := []string{lessonID}
	if !IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		fatalLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lessonMaterial := new(domain.LessonMaterial)
	if err := c.Bind(lessonMaterial); err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := usecase.CreateAuthoringLessonMaterial(c.Request(), lessonID, *lessonMaterial); err != nil {
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

func putAuthoringLessonMaterial(c echo.Context) error {
	lessonID := c.Param("id")

	ids := []string{lessonID}
	if !IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		fatalLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lessonMaterial := new(domain.LessonMaterial)
	if err := c.Bind(lessonMaterial); err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := usecase.UpdateAuthoringLessonMaterial(c.Request(), lessonID, *lessonMaterial); err != nil {
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
