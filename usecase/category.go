package usecase

import (
	"net/http"
	"strconv"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// GetCategories return categories by the subject.
func GetCategories(request *http.Request) ([]domain.Category, error) {
	queryString := request.URL.Query().Get("subjectID")
	subjectID, err := strconv.ParseInt(queryString, 10, 64)

	if err != nil {
		return nil, err
	}

	ctx := request.Context()
	// Right now, only the Japanese category exists.
	return domain.GetJapaneseCategories(ctx, subjectID)
}
