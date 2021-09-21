package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func getCategory(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	subjectIDStr := c.QueryParam("subject_id")
	if subjectIDStr == "" {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	subjectID, err := strconv.ParseInt(subjectIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	category, err := usecase.GetCategory(c.Request(), id, subjectID)
	if err != nil {
		if ok := errors.Is(err, domain.CategoryNotFound); ok {
			warnLog(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, category)

}

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
