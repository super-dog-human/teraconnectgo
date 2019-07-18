package usecase

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"google.golang.org/appengine"
)

// GetStorageObjectURLs is generate signed URL of object in GCS.
func GetStorageObjectURLs(request *http.Request, fileRequests []domain.FileRequest) (domain.SignedURLs, error) {
	ctx := appengine.NewContext(request)

	urlLength := len(fileRequests)
	urls := make([]domain.SignedURL, urlLength)

	var signedURLs domain.SignedURLs
	for i, request := range fileRequests {
		// TODO check user permission
		// TODO check file exists

		url, err := domain.GetSignedURL(ctx, request)
		if err != nil {
			return signedURLs, err
		}
		urls[i] = domain.SignedURL{request.ID, url}
	}

	signedURLs = domain.SignedURLs{SignedURLs: urls}

	return signedURLs, nil
}
