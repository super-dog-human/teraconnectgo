package usecase

import (
	"net/http"

	"github.com/SuperDogHuman/teraconnectgo/domain"
	"google.golang.org/appengine"
)

// GetAvailableGraphics for fetch graphic object from Cloud Datastore
func GetAvailableGraphics(request *http.Request) ([]domain.Graphic, error) {
	ctx := appengine.NewContext(request)

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
