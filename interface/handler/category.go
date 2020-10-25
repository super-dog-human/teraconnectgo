package handler

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/usecase"
	"github.com/labstack/echo/v4"
)

func getCategories(c echo.Context) error {
	categories := usecase.GetCategories()
	return c.JSON(http.StatusOK, categories)
}
