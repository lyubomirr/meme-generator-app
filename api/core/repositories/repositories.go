package repositories

import (
	"github.com/lyubomirr/meme-generator-app/core/entities"
)
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

type File interface {
	Save(file []byte, path string) error
	Delete(path string) error
}

//go:generate mockgen -destination=../mocks/mock_uow.go -package=mocks -mock_names=UnitOfWork=MockUnitOfWork . UnitOfWork
type UnitOfWork interface {
	GetUserRepository() User
	GetMemeRepository() Meme
	GetTemplateRepository() Template
	GetFileRepository() File
	BeginTransaction() error
	CommitTransaction() error
	RollbackTransaction() error
}

//go:generate mockgen -destination=../mocks/mock_uow_factory.go -package=mocks -mock_names=UoWFactory=MockUoWFactory . UoWFactory
type UoWFactory interface {
	Create() UnitOfWork
}