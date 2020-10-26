package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

type UserErrorCode uint

const (
	UserNotAvailable  UserErrorCode = 1
	AlreadyUserExists UserErrorCode = 2
)

func (e UserErrorCode) Error() string {
	switch e {
	case UserNotAvailable:
		return "user not available"
	case AlreadyUserExists:
		return "user is already created"
	default:
		return "unknown error"
	}
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
	ctx := request.Context()

	currentUser, err := domain.GetUserByID(ctx, id)
	if err != nil {
		return currentUser, err
	}

	return currentUser, nil
}

func CreateUser(request *http.Request, user *domain.User) error {
	ctx := request.Context()

	// not error when current user was not found.
	if _, err := domain.GetCurrentUser(request); err != nil && err != domain.UserNotFound {
		return err
	} else if err == nil {
		return AlreadyUserExists
	}

	if err := domain.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func UpdateUser(request *http.Request, user *domain.User) error {
	ctx := request.Context()

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
	ctx := request.Context()

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return err
	}

	lessons, err := domain.GetLessonsByUserID(ctx, currentUser.ID)
	if err != nil {
		return err
	}

	for _, lesson := range lessons {
		if err := deleteLessonAndRecources(ctx, lesson); err != nil {
			return err
		}
	}

	if err = domain.DeleteUser(ctx, currentUser.ID); err != nil {
		return err
	}

	return nil
}
