package handler

import (
	"encoding/json"
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/usecase"
	"github.com/labstack/echo/v4"
)

func getStorageObjects(c echo.Context) error {
	request := c.Request()
	jsonString := request.Header.Get("X-Get-Params")
	var fileRequests []domain.FileRequest

	if err := json.Unmarshal([]byte(jsonString), &fileRequests); err != nil {
		fatalLog(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if len(fileRequests) == 0 {
		errMessage := "storage object(s) params not found."
		fatalLog(errMessage)
		return c.JSON(http.StatusNotFound, errMessage)
	}

	urls, err := usecase.GetStorageObjectURLs(request, fileRequests)
	if err != nil {
		fatalLog(err)
		authErr, ok := err.(domain.AuthErrorCode)
		if ok && authErr == domain.UserNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		storageObjectErr, ok := err.(usecase.StorageObjectErrorCode)
		if ok && storageObjectErr == usecase.ObjectNotAvailable {
			return c.JSON(http.StatusForbidden, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, urls)
}
