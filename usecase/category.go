package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// GetCategory return a category by id.
func GetCategory(request *http.Request, id int64, subjectID int64) (domain.Category, error) {
	ctx := request.Context()
	return domain.GetJapaneseCategory(ctx, id, subjectID)
}

// GetCategories return categories by the subject.
func GetCategories(request *http.Request, subjectID int64) ([]domain.ShortCategory, error) {
	ctx := request.Context()
	// Right now, only the Japanese category exists.
	return domain.GetJapaneseCategories(ctx, subjectID)
}

// GetCategory return all categories.
func GetAllCategories(request *http.Request) ([]domain.ShortCategory, error) {
	ctx := request.Context()
	return domain.GetAllJapaneseCategories(ctx)
}
