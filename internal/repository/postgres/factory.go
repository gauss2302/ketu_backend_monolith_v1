package postgres

import (
	"ketu_backend_monolith_v1/internal/pkg/database"
	"ketu_backend_monolith_v1/internal/repository/interfaces"
)

type Repositories struct {
	User       interfaces.UserRepository
	Restaurant interfaces.RestaurantRepository
	Owner      interfaces.OwnerRepository
}

func NewRepositories(db *database.DB) *Repositories {
	return &Repositories{
		User:       NewUserRepository(db),
		Restaurant: NewRestaurantRepository(db),
		Owner:      NewOwnerRepository(db),
	}
}
