package handler

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/usecase"
	"github.com/labstack/echo/v4"
)

func getRawVoiceTexts(c echo.Context) error {
	id := c.Param("id")

	if voiceTexts, err := usecase.GetRawVoiceTexts(c.Request(), id); err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	} else if len(voiceTexts) == 0 {
		warnLog(err)
		return c.JSON(http.StatusNotFound, "raw voice texts not found.")
	} else {
		return c.JSON(http.StatusOK, voiceTexts)
	}
}
