package usecase

import (
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/imdario/mergo"
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

	var pendingKey *datastore.PendingKey
	commit, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		err = domain.ReserveUserProviderIDInTransaction(tx, providerID)
		if err == domain.AlreadyProviderIDExists {
			return AlreadyUserExists
		}
		if err != nil {
			return err
		}

		if pendingKey, err = domain.CreateUserInTransaction(tx, user); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	user.ID = commit.Key(pendingKey).ID

	return nil
}

func UpdateUser(request *http.Request, user *domain.User) error {
	ctx := request.Context()

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return UserNotAvailable
	}

	if err := mergo.Merge(user, currentUser); err != nil {
		return err
	}
	user.Created = currentUser.Created // Created field not merged because this time.Time fieled was initialized is not nil.

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
