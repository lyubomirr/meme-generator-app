package repositories

import "github.com/lyubomirr/meme-generator-app/core/entities"
//go:generate mockgen -destination=../mocks/mock_user_repository.go -package=mocks -mock_names=User=MockUserRepository . User
type User interface {
	Get(id uint) (entities.User, error)
	GetByUsername(username string) (entities.User, error)
	Create(user entities.User) (entities.User, error)
	Update(user entities.User) (entities.User, error)
	Delete(id uint) error
}