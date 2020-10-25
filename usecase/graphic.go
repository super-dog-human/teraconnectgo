package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/rs/xid"
	"google.golang.org/appengine"
)

// GetAvailableGraphics for fetch graphic object from Cloud Datastore
func GetAvailableGraphics(request *http.Request) ([]domain.Graphic, error) {
	ctx := appengine.NewContext(request)

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
	ctx := appengine.NewContext(request)

	var signedURLs domain.SignedURLs

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return signedURLs, err
	}

	urls := make([]domain.SignedURL, len(objectRequest.FileRequests))

	for i, fileRequest := range objectRequest.FileRequests {
		fileID := xid.New().String()

		url, err := domain.CreateBlankFileToGCS(ctx, fileID, "graphic", fileRequest)
		if err != nil {
			return signedURLs, err
		}
		urls[i] = domain.SignedURL{fileID, url}

		if err = domain.CreateGraphic(ctx, fileID, currentUser.ID, fileRequest.Extension); err != nil {
			return signedURLs, err
		}
	}

	signedURLs = domain.SignedURLs{SignedURLs: urls}
	return signedURLs, nil
}
