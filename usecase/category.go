package usecase

import "github.com/super-dog-human/teraconnectgo/domain"

// GetCategories for fetch avatar object from Cloud Datastore
func GetCategories() []domain.Category {
	return domain.GetAllCategories()
}
