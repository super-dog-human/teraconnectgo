package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func getBackgroundMusics(c echo.Context) error {
	musics, err := usecase.GetBackgroundMusics(c.Request())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, musics)
}
