package usecase

import (
	"context"
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"github.com/SuperDogHuman/teraconnectgo/infrastructure"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// GetAvailableAvatars for fetch avatar object from Cloud Datastore
func GetAvailableAvatars(request *http.Request) ([]domain.Avatar, error) {
	ctx := appengine.NewContext(request)

	var avatars []domain.Avatar

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return nil, err
	}

	usersAvatars, err := getCurrentUsersAvatars(ctx, currentUser.ID)
	if err != nil {
		return nil, err
	}
	avatars = append(avatars, usersAvatars...)

	publicAvatars, err := getPublicAvatars(ctx)
	if err != nil {
		return nil, err
	}
	avatars = append(avatars, publicAvatars...)

	return avatars, nil
}

func GetAvatarByIds(ctx context.Context, avatarID string) (domain.Avatar, error) {
	avatar := new(domain.Avatar)

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

func getCurrentUsersAvatars(ctx context.Context, userID string) ([]domain.Avatar, error){
	var avatars []domain.Avatar

	query := datastore.NewQuery("Avatar").Filter("UserId =", userID)
	keys, err := query.GetAll(ctx, &avatars)
	if err != nil {
		return nil, err
	}

	storeAvatarThumbnailUrl(ctx, &avatars, keys)

	return avatars, nil
}

func getPublicAvatars(ctx context.Context) ([]domain.Avatar, error){
	var avatars []domain.Avatar

	query := datastore.NewQuery("Avatar").Filter("IsPublic =", true)
	keys, err := query.GetAll(ctx, &avatars)
	if err != nil {
		return nil, err
	}

	storeAvatarThumbnailUrl(ctx, &avatars, keys)

	return avatars, nil
}

func storeAvatarThumbnailUrl(ctx context.Context, avatars *[]domain.Avatar, keys []*datastore.Key) {
	for i, key := range keys {
		id := key.StringID()
		(*avatars)[i].ID = id
		(*avatars)[i].ThumbnailURL = infrastructure.AvatarThumbnailURL(ctx, id)
	}
}