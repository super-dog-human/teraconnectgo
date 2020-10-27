package domain

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

func GetAvatarByIDs(ctx context.Context, id string) (Avatar, error) {
	avatar := new(Avatar)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *avatar, err
	}

	key := datastore.NameKey("Avatar", id, nil)
	if err := client.Get(ctx, key, avatar); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return *avatar, err
		}
		return *avatar, err
	}

	avatar.ID = id

	return *avatar, nil
}

func GetCurrentUsersAvatars(ctx context.Context, userID string) ([]Avatar, error) {
	var avatars []Avatar

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return nil, err
	}

	query := datastore.NewQuery("Avatar").Filter("UserID =", userID)
	keys, err := client.GetAll(ctx, query, &avatars)
	if err != nil {
		return nil, err
	}

	storeAvatarThumbnailURL(ctx, &avatars, keys)

	return avatars, nil
}

func GetPublicAvatars(ctx context.Context) ([]Avatar, error) {
	var avatars []Avatar

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return avatars, err
	}

	query := datastore.NewQuery("Avatar").Filter("IsPublic =", true)
	keys, err := client.GetAll(ctx, query, &avatars)
	if err != nil {
		return nil, err
	}

	storeAvatarThumbnailURL(ctx, &avatars, keys)

	return avatars, nil
}

func CreateAvatar(ctx context.Context, id string, userID string) error {
	avatar := new(Avatar)

	avatar.ID = id
	avatar.UserID = userID
	avatar.Created = time.Now()

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	key := datastore.NameKey("Avatar", avatar.ID, nil)
	if _, err := client.Put(ctx, key, avatar); err != nil {
		return err
	}

	return nil
}

func storeAvatarThumbnailURL(ctx context.Context, avatars *[]Avatar, keys []*datastore.Key) {
	for i, key := range keys {
		id := key.Name
		(*avatars)[i].ID = id
		(*avatars)[i].ThumbnailURL = infrastructure.AvatarThumbnailURL(ctx, id)
	}
}

func DeleteAvatarInTransaction(tx *datastore.Transaction, id string) error {
	key := datastore.NameKey("Avatar", id, nil)
	if err := tx.Delete(key); err != nil {
		return err
	}

	return nil
}
