package handler

import (
	"net/http"
	"strconv"

	"github.com/super-dog-human/teraconnectgo/domain"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func getLessonGraphics(c echo.Context) error {
	lessonID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
	}

	viewKey := c.QueryParam("view_key")
	lesson, err := usecase.GetPublicLesson(c.Request(), lessonID, viewKey)
	if err != nil {
		if err == usecase.LessonNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		} else if err == usecase.LessonNotAvailable {
			warnLog(err)
			return c.JSON(http.StatusForbidden, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ids := c.Request().URL.Query()["ids"]

	if len(ids) == 0 {
		return c.JSON(http.StatusNotFound, "invalid params.")
	}

	urls, err := usecase.GetGraphicsByLessonIDAndIDs(c.Request(), lessonID, lesson.UserID, ids)

	if err != nil {
		fatalLog(err)
		if err == domain.GraphicNotFound {
			// idsパラメータを持つ全部または一部のGraphicが見つからなかった場合
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, urls)
}
