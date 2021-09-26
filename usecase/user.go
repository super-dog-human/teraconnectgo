package usecase

import (
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
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

// GetCurrentUser for fetch current user account
func GetCurrentUser(request *http.Request) (domain.User, error) {
	var currentUser domain.User
	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return currentUser, err
	}

	return currentUser, nil
}

// GetUser for fetch user account by id.
func GetUser(request *http.Request, id int64) (domain.User, error) {
	ctx := request.Context()
	var user domain.User

	user, err := domain.GetUserByID(ctx, id)
	if err != nil {
		return user, err
	}

	user.ID = id

	return user, nil
}

// CreateUser creates new user with exclusion control.
func CreateUser(request *http.Request, user *domain.User) error {
	ctx := request.Context()

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	providerID, err := domain.ProviderID(request)
	if err != nil {
		return err
	}

	user.ProviderID = providerID

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		err = domain.ReserveUserProviderIDInTransaction(tx, providerID)
		if err == domain.AlreadyProviderIDExists {
			return AlreadyUserExists
		}
		if err != nil {
			return err
		}

		if _, err = domain.CreateUserInTransaction(tx, user); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func UpdateUser(request *http.Request, params *map[string]interface{}) (domain.User, error) {
	ctx := request.Context()

	user, err := domain.GetCurrentUser(request)
	if err != nil {
		return user, UserNotAvailable
	}

	targetFields := []string{"Name", "Profile", "Email"}
	if err = domain.UpdateUser(ctx, &user, params, &targetFields); err != nil {
		return user, err
	}

	return user, nil
}

func UnsubscribeCurrentUser(request *http.Request) error {
	ctx := request.Context()

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return err
	}

	if err = domain.DeleteUser(ctx, currentUser.ID); err != nil {
		return err
	}

	return nil
}
