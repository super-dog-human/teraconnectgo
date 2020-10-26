package usecase

import (
	"net/http"

	"github.com/rs/xid"
	"github.com/super-dog-human/teraconnectgo/domain"
)

// CreateBlankRawVoiceFile for create blank file of raw voice text to Cloud Storage.
func CreateBlankRawVoiceFile(request *http.Request, lessonID string) (domain.SignedURL, error) {
	ctx := request.Context()

	var signedURL domain.SignedURL
	if err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return signedURL, err
	}

	fileID := xid.New().String()

	if err := domain.CreateBlankRawVoiceFile(ctx, lessonID, fileID); err != nil {
		return signedURL, err
	}

	url, err := domain.GetRawVoiceFileSignedURLForUpload(ctx, lessonID, fileID)
	if err != nil {
		return signedURL, err
	}

	signedURL = domain.SignedURL{FileID: fileID, SignedURL: url}
	return signedURL, nil
}
