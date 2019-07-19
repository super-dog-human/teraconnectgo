package usecase

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"google.golang.org/appengine"
)

type UserErrorCode uint

const (
	UserNotAvailable  UserErrorCode = 1
	AlreadyUserExists UserErrorCode = 2
)

func (_ UserErrorCode) Error() string {
	return "user not available"
}

// GetUser for fetch current user account
func GetCurrentUser(request *http.Request) (domain.User, error) {
	var currentUser domain.User
	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return currentUser, err
	}

	return currentUser, nil
}

// GetUser for fetch user account by id.
func GetUser(request *http.Request, id string) (domain.User, error) {
	ctx := appengine.NewContext(request)

	currentUser, err := domain.GetUserByID(ctx, id)
	if err != nil {
		return currentUser, err
	}

	return currentUser, nil
}

func CreateUser(request *http.Request, user domain.User) error {
	ctx := appengine.NewContext(request)

	// not error when current user not found.
	if _, err := domain.GetCurrentUser(request); err != domain.UserNotFound {
		return err
	} else if err == nil {
		return AlreadyUserExists
	}

	if err := domain.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func UpdateUser(request *http.Request, user domain.User) error {
	ctx := appengine.NewContext(request)

	currentUser, err := domain.GetCurrentUser(request)
	if err == nil {
		return err
	}

	if user.ID != currentUser.ID {
		return UserNotAvailable
	}

	if err = domain.UpdateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func UnsubscribeCurrentUser(request *http.Request) error {
	ctx := appengine.NewContext(request)

	currentUser, err := domain.GetCurrentUser(request)
	if err == nil {
		return err
	}

	if err := domain.DestroyUser(ctx, currentUser.ID); err != nil {
		return err
	}

	// remove all resources with user.

	return err
}