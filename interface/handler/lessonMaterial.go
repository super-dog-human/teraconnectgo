package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func getLessonMaterials(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lessonID, err := strconv.ParseInt(c.Param("lessonID"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lessonMaterial, err := usecase.GetLessonMaterial(c.Request(), id, lessonID)
	if err != nil {
		fatalLog(err)
		lessonErr, ok := err.(usecase.LessonMaterialErrorCode)
		if ok && lessonErr == usecase.LessonMaterialNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok && lessonErr == usecase.LessonMaterialNotAvailable {
			return c.JSON(http.StatusForbidden, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lessonMaterial)
}

func postLessonMaterial(c echo.Context) error {
	lessonID, err := strconv.ParseInt(c.Param("lessonID"), 10, 64)
	if err != nil {
		errMessage := "Invalid lessonID error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	params := new(usecase.LessonMaterialParams)
	if err := c.Bind(params); err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	id, err := usecase.CreateLessonMaterial(c.Request(), lessonID, *params)
	if err != nil {
		fatalLog(err)
		lessonErr, ok := err.(usecase.LessonMaterialErrorCode)
		if ok && lessonErr == usecase.LessonMaterialNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok && lessonErr == usecase.LessonMaterialNotAvailable {
			return c.JSON(http.StatusForbidden, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := struct {
		ID int64 `json:"materialID"`
	}{
		id,
	}

	return c.JSON(http.StatusCreated, response)
}

func patchLessonMaterial(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	lessonID, err := strconv.ParseInt(c.Param("lessonID"), 10, 64)
	if err != nil {
		errMessage := "Invalid lessonID error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	params := new(usecase.LessonMaterialParams)
	if err := c.Bind(params); err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := usecase.UpdateLessonMaterial(c.Request(), id, lessonID, *params); err != nil {
		fatalLog(err)
		lessonErr, ok := err.(usecase.LessonMaterialErrorCode)
		if ok && lessonErr == usecase.LessonMaterialNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if ok && lessonErr == usecase.LessonMaterialNotAvailable {
			return c.JSON(http.StatusForbidden, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "succeeded")
}
