package handler

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func getBackgroundMusics(c echo.Context) error {
	musics, err := usecase.GetBackgroundMusics(c.Request())
	if err != nil {
		_, ok := err.(domain.AuthErrorCode)
		if ok {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		fatalLog(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, musics)
}

func postBackgroundMusic(c echo.Context) error {
	param := new(usecase.CreateBackgroundMusicParam)

	if err := c.Bind(param); err != nil {
		fatalLog(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	param.Name = strings.TrimSpace(param.Name)
	if len(param.Name) == 0 {
		return c.JSON(http.StatusBadRequest, "invalid name.")
	}

	if utf8.RuneCountInString(param.Name) > 50 {
		param.Name = string([]rune(param.Name)[:50])
	}

	signedURL, err := usecase.CreateBackgroundMusicAndBlankFile(c.Request(), param)
	if err != nil {
		_, ok := err.(domain.AuthErrorCode)
		if ok {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, signedURL)
}
