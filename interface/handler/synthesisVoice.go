package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func postSynthesisVoice(c echo.Context) error {
	param := new(domain.CreateSynthesisVoiceParam)

	if err := c.Bind(param); err != nil {
		fatalLog(err)
		authErr, ok := err.(domain.AuthErrorCode)
		if ok && authErr == domain.UserNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	voice, err := usecase.CreateSynthesisVoice(c.Request(), param)
	if err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := voiceResponse{ID: voice.ID, FileKey: voice.FileKey}

	return c.JSON(http.StatusOK, response)
}

type voiceResponse struct {
	ID      int64  `json:"id"`
	FileKey string `json:"fileKey"`
}
