package usecase

import (
	"context"
	"net/http"
	"strings"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"google.golang.org/appengine"
)


type StorageObjectErrorCode uint

const ObjectNotAvailable StorageObjectErrorCode = 1

func (_ StorageObjectErrorCode) Error() string {
	return "object not available"
}

// GetStorageObjectURLs is generate signed URL of object in GCS.
func GetStorageObjectURLs(request *http.Request, fileRequests []domain.FileRequest) (domain.SignedURLs, error) {
	ctx := appengine.NewContext(request)

	var signedURLs domain.SignedURLs

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return signedURLs, err
	}

	urlLength := len(fileRequests)
	urls := make([]domain.SignedURL, urlLength)

	for i, fileRequest := range fileRequests {
		if err = currentUserAccessToStorageObject(ctx, request, fileRequest, currentUser.ID); err != nil {
			return signedURLs, err
		}

		url, err := domain.GetSignedURL(ctx, fileRequest)
		if err != nil {
			return signedURLs, err
		}
		urls[i] = domain.SignedURL{fileRequest.ID, url}
	}

	signedURLs = domain.SignedURLs{SignedURLs: urls}

	return signedURLs, nil
}

func currentUserAccessToStorageObject(ctx context.Context, request *http.Request, fileRequest domain.FileRequest, userID string) error {
	rawEntityName := strings.ToLower(fileRequest.Entity)
	entityID, entityName := entityIDFromRequest(rawEntityName, fileRequest.ID)
	entity, err := domain.EntityOfRequestedFile(ctx, entityID, entityName)
	if err != nil {
		return err
	}

	if entity.UserID != userID {
		return ObjectNotAvailable
	}

	return nil
}

func entityIDFromRequest(entityName string, rawID string) (string, string) {
	switch entityName {
		case "Lesson":
			return rawID, entityName
		case "Avatar":
			return rawID, entityName
		case "Graphic":
			return rawID, entityName
		default:
			// using Lesson when entity is "voice/:lessonID"
			separatorIndex := strings.Index(rawID, "/")
			return rawID[separatorIndex:len(rawID)], "Lesson"
	}
}
