package usecase

import (
	"net/http"
	"strconv"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// GetAvailableAvatars for fetch avatar object from Cloud Datastore
func GetAvailableAvatars(request *http.Request) ([]domain.Avatar, error) {
	ctx := request.Context()

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
	ctx := request.Context()

	var signedURLs domain.SignedURLs

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return signedURLs, err
	}

	urls := make([]domain.SignedURL, len(objectRequest.FileRequests))

	for i, fileRequest := range objectRequest.FileRequests {
		avatar := new(domain.Avatar)

		if err = domain.CreateAvatar(ctx, avatar, &currentUser); err != nil {
			return signedURLs, err
		}

		fileID := strconv.FormatInt(avatar.ID, 10)
		url, err := domain.CreateBlankFileToGCS(ctx, fileID, "avatar", fileRequest)
		if err != nil {
			return signedURLs, err
		}
		urls[i] = domain.SignedURL{FileID: fileID, SignedURL: url}

	}

	signedURLs = domain.SignedURLs{SignedURLs: urls}
	return signedURLs, nil
}
