package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func patchLessonViewCount(c echo.Context) error {
	lessonID, err := strconv.ParseInt(c.QueryParam("lesson_id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	if err = usecase.UpdateLessonViewCount(c.Request(), lessonID); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "succeeded")
}
