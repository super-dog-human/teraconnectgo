package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// GetBackgroundMusics returns music URLs in Cloud Datastore.
func GetBackgroundMusics(request *http.Request) ([]domain.BackgroundMusic, error) {
	ctx := request.Context()
	return domain.GetAllBackgroundMusics(ctx)
}
