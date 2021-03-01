package repositories

import "github.com/lyubomirr/meme-generator-app/core/entities"
//go:generate mockgen -destination=../mocks/mock_user_repository.go -package=mocks -mock_names=User=MockUserRepository . User
type User interface {
	Get(id uint) (entities.User, error)
	GetByUsername(username string) (entities.User, error)
	Create(user entities.User) (uint, error)
	Update(user entities.User) (entities.User, error)
	Delete(id uint) error
}

type Meme interface {
	GetAll() ([]entities.Meme, error)
	Get(id uint) (entities.Meme, error)
	GetByAuthor(userId uint) ([]entities.Meme, error)
	Create(meme entities.Meme) (uint, error)
	Update(meme entities.Meme) (entities.Meme, error)
	Delete(id uint) error
}

type Template interface {
	GetAll() ([]entities.Template, error)
	Get(id uint) (entities.Template, error)
	Create(meme entities.Template) (uint, error)
	Update(meme entities.Template) (entities.Template, error)
	Delete(id uint) error
}