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

func CreateGraphicsAndBlankFiles(request *http.Request, objectRequest domain.StorageObjectRequest) (domain.SignedURLs, error) {
	ctx := request.Context()

	var signedURLs domain.SignedURLs

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return signedURLs, err
	}

	lesson, err := domain.GetLessonByID(ctx, objectRequest.LesonID)
	if err != nil {
		return signedURLs, err
	}

	if lesson.UserID != currentUser.ID {
		return signedURLs, LessonNotAvailable
	}

	graphics := make([]*domain.Graphic, len(objectRequest.FileRequests))
	urls := make([]domain.SignedURL, len(objectRequest.FileRequests))

	for i, fileRequest := range objectRequest.FileRequests {
		graphic := new(domain.Graphic)
		graphic.LessonID = objectRequest.LesonID
		graphic.FileType = fileRequest.Extension
		graphics[i] = graphic
	}

	if err = domain.CreateGraphics(ctx, &currentUser, graphics); err != nil {
		return signedURLs, err
	}

	for i, fileRequest := range objectRequest.FileRequests {
		fileID := strconv.FormatInt(graphics[i].ID, 10)
		url, err := domain.CreateBlankFileToGCS(ctx, fileID, "graphic", fileRequest)
		if err != nil {
			return signedURLs, err
		}
		urls[i] = domain.SignedURL{FileID: fileID, SignedURL: url}
	}

	return domain.SignedURLs{SignedURLs: urls}, nil
}
