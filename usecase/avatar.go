package usecase

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"google.golang.org/appengine"
)

// GetAvailableAvatars for fetch avatar object from Cloud Datastore
func GetAvailableAvatars(request *http.Request) ([]domain.Avatar, error) {
	ctx := appengine.NewContext(request)

	var avatars []domain.Avatar

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return nil, err
	}

	usersAvatars, err := domain.GetCurrentUsersAvatars(ctx, currentUser.ID)
	if err != nil {
		return nil, err
	}
	avatars = append(avatars, usersAvatars...)

	publicAvatars, err := domain.GetPublicAvatars(ctx)
	if err != nil {
		return nil, err
	}
	avatars = append(avatars, publicAvatars...)

	return avatars, nil
}

