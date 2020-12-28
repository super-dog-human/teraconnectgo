package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// GetBackgroundImages returns image URLs in Cloud Datastore.
func GetBackgroundImages(request *http.Request) ([]domain.BackgroundImage, error) {
	ctx := request.Context()
	return domain.GetAllBackgroundImages(ctx)
}
