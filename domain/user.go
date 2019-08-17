package domain

import (
	"context"
	"time"

	"google.golang.org/appengine/datastore"
	"github.com/rs/xid"
)

// GetUserByID is return user has ID.
func GetUserByID(ctx context.Context, id string) (User, error) {
	user := new(User)

	key := datastore.NewKey(ctx, "User", id, 0, nil)
	if err := datastore.Get(ctx, key, user); err != nil {
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
	user.Created = time.Now()

	user.ID = xid.New().String()

	key := datastore.NewKey(ctx, "User", user.ID, 0, nil)
	if _, err := datastore.Put(ctx, key, user); err != nil {
		return err
	}

	return nil
}

func UpdateUser(ctx context.Context, user *User) error {
	user.Updated = time.Now()

	key := datastore.NewKey(ctx, "User", user.ID, 0, nil)
	if _, err := datastore.Put(ctx, key, user); err != nil {
		return err
	}

	return nil
}

func DeleteUser(ctx context.Context, id string) error {
	key := datastore.NewKey(ctx, "User", id, 0, nil)
	if err := datastore.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}