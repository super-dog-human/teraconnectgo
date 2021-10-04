package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

type getLessonMaterialShortResponse struct {
	DurationSec          float32                     `json:"durationSec"`
	AvatarID             int64                       `json:"avatarID"`
	Avatar               domain.Avatar               `json:"avatar"`
	AvatarLightColor     string                      `json:"avatarLightColor"`
	BackgroundImageID    int64                       `json:"backgroundImageID"`
	BackgroundImageURL   string                      `json:"backgroundImageURL"`
	VoiceSynthesisConfig domain.VoiceSynthesisConfig `json:"voiceSynthesisConfig"`
	Created              time.Time                   `json:"created"`
	Updated              time.Time                   `json:"updated"`
}

type postMaterialResponse struct {
	MaterialID int64 `json:"materialID"`
}

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

	isShort := c.Request().URL.Query().Get("is_short")
	if isShort == "true" {
		var response getLessonMaterialShortResponse
		copier.Copy(&response, lessonMaterial)
		return c.JSON(http.StatusOK, response)

	} else {
		return c.JSON(http.StatusOK, lessonMaterial)
	}
}

func patchLessonMaterial(c echo.Context) error {
	lessonID, err := strconv.ParseInt(c.Param("lessonID"), 10, 64)
	if err != nil {
		errMessage := "Invalid lessonID error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	var params map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := usecase.UpdateLessonMaterial(c.Request(), id, lessonID, &params); err != nil {
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
