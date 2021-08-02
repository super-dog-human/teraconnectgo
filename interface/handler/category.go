package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func getCategories(c echo.Context) error {
	queryString := c.QueryParam("subject_id")
	subjectID, err := strconv.ParseInt(queryString, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	categories, err := usecase.GetCategories(c.Request(), subjectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, categories)

}
