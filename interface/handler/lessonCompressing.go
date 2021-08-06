package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

func postLessonCompressing(c echo.Context) error {
	request := c.Request()

	queueName := request.Header.Get("X-Appengine-Queuename")
	if queueName != "compressLesson" {
		log.Printf("invalid queueName: %v\n", queueName)
		return c.JSON(http.StatusBadRequest, "")
	}

	taskName := request.Header.Get("X-Appengine-Taskname")
	if taskName == "" {
		log.Printf("invalid taskName: %v\n", taskName)
		return c.JSON(http.StatusBadRequest, "")
	}

	tasks := strings.Split(taskName, "-")
	lessonID, err := strconv.ParseInt(tasks[1], 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	queuedUnixNanoTime, err := strconv.ParseInt(tasks[2], 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err = usecase.CompreesLesson(request, lessonID, taskName, queuedUnixNanoTime); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "succeeded.")
}
