package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// GetCategories return categories by the subject.
func GetCategories(request *http.Request, subjectID int64) ([]domain.Category, error) {
	ctx := request.Context()
	// Right now, only the Japanese category exists.
	return domain.GetJapaneseCategories(ctx, subjectID)
}
