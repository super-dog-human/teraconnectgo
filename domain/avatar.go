package domain

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// GetAvatarByIDs gets avatar by id.
func GetAvatarByIDs(ctx context.Context, id int64) (Avatar, error) {
	avatar := new(Avatar)

	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return *avatar, err
	}

	key := datastore.IDKey("Avatar", id, nil)
	if err := client.Get(ctx, key, avatar); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return *avatar, err
		}
		return *avatar, err
	}

	avatar.ID = id

	return *avatar, nil
}

// GetCurrentUsersAvatars gets avatars belongs to user.
func GetCurrentUsersAvatars(ctx context.Context, userID int64) ([]Avatar, error) {
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

func CreateAvatar(ctx context.Context, avatar *Avatar) error {
	client, err := datastore.NewClient(ctx, infrastructure.ProjectID())
	if err != nil {
		return err
	}

	avatar.Created = time.Now()

	key := datastore.IncompleteKey("Avatar", nil)
	putKey, err := client.Put(ctx, key, avatar)
	if err != nil {
		return err
	}

	avatar.ID = putKey.ID

	return nil
}

func storeAvatarThumbnailURL(ctx context.Context, avatars *[]Avatar, keys []*datastore.Key) {
	for i, key := range keys {
		(*avatars)[i].ID = key.ID
		(*avatars)[i].ThumbnailURL = infrastructure.AvatarThumbnailURL(ctx, key.ID)
	}
}

func DeleteAvatarInTransaction(tx *datastore.Transaction, id int64) error {
	key := datastore.IDKey("Avatar", id, nil)
	if err := tx.Delete(key); err != nil {
		return err
	}

	return nil
}
