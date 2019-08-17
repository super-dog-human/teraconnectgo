package handler

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/SuperDogHuman/teraconnectgo/usecase"
	"github.com/labstack/echo/v4"
)

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
	id := c.Param("id")

	ids := []string{id}
	if !IsValidXIDs(ids) {
		errMessage := "Invalid ID(s) error"
		warnLog(errMessage)
		return c.JSON(http.StatusBadRequest, errMessage)
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

func postUser(c echo.Context) error {
	user := new(domain.User)

	userSubject, err := domain.UserSubject(c.Request())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user.Auth0Sub = userSubject

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
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
	user := new(domain.User)

	if err := c.Bind(user); err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := usecase.UpdateUser(c.Request(), user); err != nil {
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
