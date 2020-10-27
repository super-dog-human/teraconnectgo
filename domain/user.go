package domain

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/rs/xid"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

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

// CreateUser is creating user.
func CreateUser(ctx context.Context, user *User) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	user.Created = time.Now()
	user.ID = xid.New().String()

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
