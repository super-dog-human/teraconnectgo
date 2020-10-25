package usecase

import "github.com/super-dog-human/teraconnectgo/domain"

// GetAvailableAvatars for fetch avatar object from Cloud Datastore
func GetCategories() []domain.Category {
	return domain.GetAllCategories()
}
