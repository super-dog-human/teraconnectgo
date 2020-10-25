package handler

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/usecase"
	"github.com/labstack/echo/v4"
)

func getRawVoiceTexts(c echo.Context) error {
	id := c.Param("id")

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
