package handler

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/SuperDogHuman/teraconnectgo/usecase"
	"github.com/labstack/echo/v4"
)

func getGraphics(c echo.Context) error {
	// TODO pagination.
	graphics, err := usecase.GetAvailableGraphics(c.Request())

	if err != nil {
		authErr, ok := err.(domain.AuthErrorCode)
		if ok && authErr == domain.UserNotFound {
			// when token is valid but user account not exists.
			return c.JSON(http.StatusNotFound, authErr)
		} else {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	if len(graphics) == 0 {
		errMessage := "graphics not found"
		warnLog(errMessage)
		return c.JSON(http.StatusNotFound, errMessage)
	}

	return c.JSON(http.StatusOK, graphics)
}
