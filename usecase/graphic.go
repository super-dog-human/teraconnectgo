package usecase

import (
	"errors"
	"net/http"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/infrastructure"
)

// GetGraphicByID is fetching a graphic by id.
func GetGraphicByID(request *http.Request, id int64) (domain.Graphic, error) {
	ctx := request.Context()

	var graphic domain.Graphic

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return graphic, err
	}

	graphic, err = domain.GetGraphicByID(ctx, id, currentUser.ID)
	if err != nil {
		return graphic, err
	}

	url, err := domain.GetGraphicSignedURL(ctx, &graphic)
	if err != nil {
		return graphic, err
	}

	graphic.URL = url

	return graphic, nil
}

// GetGraphicsByLessonID is fetching graphics belongs to lesson.
func GetGraphicsByLessonID(request *http.Request, lessonID int64) ([]*domain.Graphic, error) {
	ctx := request.Context()

	var graphics []*domain.Graphic

	if _, err := currentUserAccessToLesson(ctx, request, lessonID); err != nil {
		return nil, err
	}

	if err := domain.GetGraphicsByLessonID(ctx, lessonID, &graphics); err != nil {
		return nil, err
	}

	return graphics, nil
}

func GetGraphicsByLessonIDAndIDs(request *http.Request, lessonID int64, userID int64, ids []string) (map[int64]string, error) {
	ctx := request.Context()

	intIDs := make([]int64, len(ids))
	for i, id := range ids {
		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		intIDs[i] = intID
	}

	graphics, err := domain.GetGraphicsByIDs(ctx, userID, intIDs)
	if err != nil {
		return nil, err
	}

	urls := make(map[int64]string)
	for _, graphic := range graphics {
		if graphic.LessonID != lessonID {
			continue // GraphicのIDさえ分かれば今回取得分のLessonに無関係なものも取得できてしまうので、紐付きを確認する
		}
		url, err := domain.GetGraphicSignedURL(ctx, graphic)
		if err != nil {
			return nil, err
		}

		urls[graphic.ID] = url
	}

	return urls, nil
}

func CreateGraphicsAndBlankFiles(request *http.Request, objectRequest infrastructure.StorageObjectRequest) (infrastructure.SignedURLs, error) {
	ctx := request.Context()

	var signedURLs infrastructure.SignedURLs

	userID, err := currentUserAccessToLesson(ctx, request, objectRequest.LessonID)
	if err != nil {
		return signedURLs, err
	}

	graphics := make([]*domain.Graphic, len(objectRequest.FileRequests))
	urls := make([]infrastructure.SignedURL, len(objectRequest.FileRequests))

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
		url, err := infrastructure.CreateBlankFileToGCS(ctx, fileID, "graphic", fileRequest)
		if err != nil {
			return signedURLs, err
		}
		urls[i] = infrastructure.SignedURL{FileID: fileID, SignedURL: url}
	}

	return infrastructure.SignedURLs{SignedURLs: urls}, nil
}

func DeleteGraphic(request *http.Request, id int64) error {
	ctx := request.Context()

	currentUser, err := domain.GetCurrentUser(request)
	if err != nil {
		return err
	}

	graphic, err := domain.GetGraphicByID(ctx, id, currentUser.ID)
	if err != nil {
		return err
	}

	if err := domain.DeleteGraphicByID(ctx, graphic.ID, currentUser.ID); err != nil {
		return err
	}

	if err := domain.DeleteGraphicFileByID(ctx, graphic); err != nil {
		if ok := errors.Is(err, storage.ErrObjectNotExist); ok {
			return nil // 削除しようとするファイルが存在しなくてもエラーにしない
		}
		return err
	}

	return nil
}
