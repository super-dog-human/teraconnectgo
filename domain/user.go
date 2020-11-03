package domain

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/rs/xid"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

type UserErrorCode uint

const (
	AlreadyProviderIDExists UserErrorCode = 1
)

func (e UserErrorCode) Error() string {
	switch e {
	case AlreadyProviderIDExists:
		return "provider id is already existed"
	default:
		return "unknown error"
	}
}

// GetCurrentUser is return logged in user
func GetCurrentUser(request *http.Request) (User, error) {
	user := new(User) // for return blank user when error

	providerID, err := ProviderID(request)
	if err != nil {
		return *user, err
	}

	var users []User
	ctx := request.Context()
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *user, FailedDatastoreInitialize
	}

	query := datastore.NewQuery("User").Filter("ProviderID =", providerID).Limit(1)
	keys, err := client.GetAll(ctx, query, &users)
	if err != nil {
		return *user, FailedGettingUser
	}

	if len(users) == 0 {
		return *user, UserNotFound
	}

	user = &users[0]
	user.ID = keys[0].Name
	return users[0], nil
}

// GetUserByID is return user has ID.
func GetUserByID(ctx context.Context, id string) (User, error) {
	user := new(User)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *user, err
	}

	key := datastore.NameKey("User", id, nil)
	if err := client.Get(ctx, key, user); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return *user, UserNotFound
		}
		return *user, err
	}
	user.ID = id

	return *user, nil
}

// ReserveUserProviderID creates user's ProviderID for  exclusion control.
func ReserveUserProviderID(request *http.Request, providerID string) error {
	ctx := request.Context()

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		key := datastore.NameKey("UserProviderID", providerID, nil)
		userProviderID := new(UserProviderID)

		err := tx.Get(key, userProviderID)
		if err == nil {
			return AlreadyProviderIDExists
		}
		if err != datastore.ErrNoSuchEntity {
			return err
		}

		// Put only when ErrNoSuchEntity
		if _, err := tx.Put(key, userProviderID); err != nil {
			return err
		}

		return nil
	})

	return err
}

// CreateUser is creating user.
func CreateUser(request *http.Request, user *User) error {
	ctx := request.Context()

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	user.ID = xid.New().String()
	user.Created = time.Now()

	key := datastore.NameKey("User", user.ID, nil)
	if _, err := client.Put(ctx, key, user); err != nil {
		return err
	}

	return nil
}

func UpdateUser(ctx context.Context, user *User) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	user.Updated = time.Now()

	key := datastore.NameKey("User", user.ID, nil)
	if _, err := client.Put(ctx, key, user); err != nil {
		return err
	}

	return nil
}

func DeleteUser(ctx context.Context, id string) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	key := datastore.NameKey("User", id, nil)
	if err := client.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}
