package usecase

import "github.com/SuperDogHuman/teraconnectgo/domain"

// GetAvailableAvatars for fetch avatar object from Cloud Datastore
func GetCategories() []domain.Category {
	return domain.GetAllCategories()
}
