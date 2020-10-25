package domain

import (
	"context"
	"time"

	"github.com/super-dog-human/teraconnectgo/infrastructure"
	"google.golang.org/appengine/datastore"
)

func GetAvatarByIds(ctx context.Context, avatarID string) (Avatar, error) {
	avatar := new(Avatar)

	avatarKey := datastore.NewKey(ctx, "Avatar", avatarID, 0, nil)
	if err := datastore.Get(ctx, avatarKey, avatar); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return *avatar, err
		}
		return *avatar, err
	}

	avatar.ID = avatarID
	return *avatar, nil
}

func GetCurrentUsersAvatars(ctx context.Context, userID string) ([]Avatar, error){
	var avatars []Avatar

	query := datastore.NewQuery("Avatar").Filter("UserId =", userID)
	keys, err := query.GetAll(ctx, &avatars)
	if err != nil {
		return nil, err
	}

	storeAvatarThumbnailUrl(ctx, &avatars, keys)

	return avatars, nil
}

func GetPublicAvatars(ctx context.Context) ([]Avatar, error){
	var avatars []Avatar

	query := datastore.NewQuery("Avatar").Filter("IsPublic =", true)
	keys, err := query.GetAll(ctx, &avatars)
	if err != nil {
		return nil, err
	}

	storeAvatarThumbnailUrl(ctx, &avatars, keys)

	return avatars, nil
}

func CreateAvatar(ctx context.Context, id string, userID string) error {
	avatar := new(Avatar)

	avatar.ID = id
	avatar.UserID  = userID
	avatar.Created = time.Now()

	key := datastore.NewKey(ctx, "Avatar", avatar.ID, 0, nil)
	if _, err := datastore.Put(ctx, key, avatar); err != nil {
		return err
	}

	return nil
}

func DeleteAvatar(ctx context.Context, id string) error {
	key := datastore.NewKey(ctx, "Avatar", id, 0, nil)
	if err := datastore.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}

func storeAvatarThumbnailUrl(ctx context.Context, avatars *[]Avatar, keys []*datastore.Key) {
	for i, key := range keys {
		id := key.StringID()
		(*avatars)[i].ID = id
		(*avatars)[i].ThumbnailURL = infrastructure.AvatarThumbnailURL(ctx, id)
	}
}
