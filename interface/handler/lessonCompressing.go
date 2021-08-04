package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func postLessonCompressing(c echo.Context) error {
	request := c.Request()
	queueName := request.Header.Get("X-Appengine-Queuename")
	fmt.Printf("queueName: %v", queueName)
	return c.JSON(http.StatusOK, "ok.")
}
