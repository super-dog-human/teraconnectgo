package usecase

import (
	"net/http"

	"github.com/super-dog-human/teraconnectgo/domain"
)

// GetSubjects for fetch avatar object from Cloud Datastore
func GetSubjects(request *http.Request) ([]domain.Subject, error) {
	ctx := request.Context()
	return domain.GetAllSubjects(ctx)
}
