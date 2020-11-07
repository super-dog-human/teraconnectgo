package usecase

import (
	"net/http"
	"strconv"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// GetAvailableGraphics for fetch graphic object from Cloud Datastore
func GetAvailableGraphics(request *http.Request) ([]domain.Graphic, error) {
	ctx := request.Context()

	var graphics []domain.Graphic

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return nil, err
	}

	usersGraphics, err := domain.GetCurrentUsersGraphics(ctx, currentUser.ID)
	if err != nil {
		return nil, err
	}
	graphics = append(graphics, usersGraphics...)

	publicGraphics, err := domain.GetPublicGraphics(ctx)
	if err != nil {
		return nil, err
	}
	graphics = append(graphics, publicGraphics...)

	return graphics, nil
}

func CreateGraphicsAndBlankFile(request *http.Request, objectRequest domain.StorageObjectRequest) (domain.SignedURLs, error) {
	ctx := request.Context()

	var signedURLs domain.SignedURLs

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return signedURLs, err
	}

	urls := make([]domain.SignedURL, len(objectRequest.FileRequests))

	for i, fileRequest := range objectRequest.FileRequests {
		graphic := new(domain.Graphic)
		graphic.UserID = currentUser.ID
		graphic.FileType = fileRequest.Extension

		if err = domain.CreateGraphic(ctx, graphic); err != nil {
			return signedURLs, err
		}

		fileID := strconv.FormatInt(graphic.ID, 10)
		url, err := domain.CreateBlankFileToGCS(ctx, fileID, "graphic", fileRequest)
		if err != nil {
			return signedURLs, err
		}
		urls[i] = domain.SignedURL{FileID: fileID, SignedURL: url}
	}

	signedURLs = domain.SignedURLs{SignedURLs: urls}
	return signedURLs, nil
}
