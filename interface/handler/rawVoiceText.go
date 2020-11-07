package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func getRawVoiceTexts(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	voiceTexts, err := usecase.GetRawVoiceTexts(c.Request(), id)
	if err != nil {
		fatalLog(err)
		lessonErr, ok := err.(usecase.LessonErrorCode)
		if ok && lessonErr == usecase.LessonNotAvailable {
			return c.JSON(http.StatusForbidden, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(voiceTexts) == 0 {
		errMessage := "raw voice texts not found."
		warnLog(errMessage)
		return c.JSON(http.StatusNotFound, errMessage)
	}

	return c.JSON(http.StatusOK, voiceTexts)
}
