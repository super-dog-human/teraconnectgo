package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

type getUserShortResponse struct {
	NextCursor string        `json:"nextCursor"`
	Users      []domain.User `json:"users"`
}

func getUserMe(c echo.Context) error {
	user, err := usecase.GetCurrentUser(c.Request())

	if err != nil {
		authErr, ok := err.(domain.AuthErrorCode)
		if ok && authErr == domain.UserNotFound {
			// when token is valid but user account not exists.
			warnLog(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func getUser(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	user, err := usecase.GetUser(c.Request(), id)

	if err != nil {
		authErr, ok := err.(domain.AuthErrorCode)
		if ok && authErr == domain.UserNotFound {
			warnLog(err)
			return c.JSON(http.StatusNotFound, err.Error())
		}
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func getUsers(c echo.Context) error {
	cursorStr := c.QueryParam("next_cursor")
	users, nextCursorStr, err := usecase.GetUsers(c.Request(), cursorStr)

	if err != nil {
		if ok := errors.Is(err, usecase.UserNotFound); ok {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := getUserShortResponse{Users: users, NextCursor: nextCursorStr}
	return c.JSON(http.StatusOK, response)
}

func postUser(c echo.Context) error {
	user := new(usecase.NewUserParams)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if user.Name == "" || user.Email == "" {
		return c.JSON(http.StatusBadRequest, "Invalid params")
	}

	if err := usecase.CreateUser(c.Request(), user); err != nil {
		userErr, ok := err.(usecase.UserErrorCode)
		if ok && userErr == usecase.AlreadyUserExists {
			return c.JSON(http.StatusConflict, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func patchUser(c echo.Context) error {
	var params map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := usecase.UpdateUser(c.Request(), &params)
	if err != nil {
		fatalLog(err)
		userErr, ok := err.(usecase.UserErrorCode)
		if ok && userErr == usecase.UserNotAvailable {
			return c.JSON(http.StatusForbidden, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func deleteUser(c echo.Context) error {
	if err := usecase.UnsubscribeCurrentUser(c.Request()); err != nil {
		fatalLog(err)
		authErr, ok := err.(domain.AuthErrorCode)
		if ok && authErr == domain.UserNotFound {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNoContent, "succeeded")
}
