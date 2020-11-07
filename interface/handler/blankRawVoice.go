package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func postBlankRawVoice(c echo.Context) error {
	request := new(postRawVoiceRequest)
	if err := c.Bind(request); err != nil {
		fatalLog(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	signedURL, err := usecase.CreateBlankRawVoiceFile(c.Request(), request.LessonID)
	if err != nil {
		fatalLog(err)
		lessonErr, ok := err.(usecase.LessonErrorCode)
		if ok && lessonErr == usecase.LessonNotAvailable {
			return c.JSON(http.StatusForbidden, err.Error())
		}
		authErr, ok := err.(domain.AuthErrorCode)
		if ok && authErr == domain.UserNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, signedURL)
}

type postRawVoiceRequest struct {
	LessonID int64 `json:"lesson_id"`
}
