package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func getVoices(c echo.Context) error {
	lessonID, err := strconv.ParseInt(c.QueryParam("lesson_id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	voices, err := usecase.GetVoices(c.Request(), lessonID)
	if err != nil {
		voiceErr, ok := err.(domain.VoiceErrorCode)
		if ok && voiceErr == domain.VoiceNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, voices)
}

func postVoice(c echo.Context) error {
	param := new(usecase.CreateVoiceParam)
	if err := c.Bind(param); err != nil {
		fatalLog(err)
		authErr, ok := err.(domain.AuthErrorCode)
		if ok && authErr == domain.UserNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	signedURL, err := usecase.CreateVoiceAndBlankFile(c.Request(), param)
	if err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, signedURL)
}
