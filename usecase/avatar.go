package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/rs/xid"
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

func CreateAvatarsAndBlankFile(request *http.Request, objectRequest domain.StorageObjectRequest) (domain.SignedURLs, error) {
	ctx := appengine.NewContext(request)

	var signedURLs domain.SignedURLs

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return signedURLs, err
	}

	urls := make([]domain.SignedURL, len(objectRequest.FileRequests))

	for i, fileRequest := range objectRequest.FileRequests {
		fileID := xid.New().String()

		url, err := domain.CreateBlankFileToGCS(ctx, fileID, "avatar", fileRequest)
		if err != nil {
			return signedURLs, err
		}
		urls[i] = domain.SignedURL{fileID, url}

		if err = domain.CreateAvatar(ctx, fileID, currentUser.ID); err != nil {
			return signedURLs, err
		}
	}

	signedURLs = domain.SignedURLs{SignedURLs: urls}
	return signedURLs, nil
}
