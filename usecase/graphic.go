package usecase

import (
	"net/http"
	"strconv"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// GetGraphicsByLessonID is fetching graphics belongs to lesson.
func GetGraphicsByLessonID(request *http.Request, lessonID int64) ([]domain.Graphic, error) {
	ctx := request.Context()

	var graphics []domain.Graphic

	if _, err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return nil, err
	}

	if err := domain.GetGraphicsByLessonID(ctx, lessonID, &graphics); err != nil {
		return nil, err
	}

	return graphics, nil
}

func CreateGraphicsAndBlankFiles(request *http.Request, objectRequest domain.StorageObjectRequest) (domain.SignedURLs, error) {
	ctx := request.Context()

	var signedURLs domain.SignedURLs

	userID, err := currentUserAccessToLesson(ctx, request, objectRequest.LessonID)
	if err != nil {
		return signedURLs, err
	}

	graphics := make([]*domain.Graphic, len(objectRequest.FileRequests))
	urls := make([]domain.SignedURL, len(objectRequest.FileRequests))

	for i, fileRequest := range objectRequest.FileRequests {
		graphic := new(domain.Graphic)
		graphic.LessonID = objectRequest.LessonID
		graphic.FileType = fileRequest.Extension
		graphics[i] = graphic
	}

	if err = domain.CreateGraphics(ctx, userID, graphics); err != nil {
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
