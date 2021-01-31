package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func getVoices(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
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
