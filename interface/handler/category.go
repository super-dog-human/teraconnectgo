package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func getCategories(c echo.Context) error {
	categories, err := usecase.GetCategories(c.Request())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, categories)

}
