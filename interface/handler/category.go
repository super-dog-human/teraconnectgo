package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func getCategories(c echo.Context) error {
	categories := usecase.GetCategories()
	return c.JSON(http.StatusOK, categories)
}
