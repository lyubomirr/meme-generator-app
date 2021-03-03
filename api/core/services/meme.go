package services

import (
	"context"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	"github.com/lyubomirr/meme-generator-app/core/repositories"
)

type Meme interface {
	GetAll(ctx context.Context) ([]entities.Meme, error)
	Get(ctx context.Context, id uint) (entities.Meme, error)
	GetByAuthor(ctx context.Context, userID uint) ([]entities.Meme, error)
	Create(ctx context.Context, file []byte, meme entities.Meme) (entities.Meme, error)
	AddComment(ctx context.Context, memeID uint, comment entities.Comment) (entities.Meme, error)
	DeleteComment(ctx context.Context, memeID uint, commentId uint) (entities.Meme, error)
	Delete(ctx context.Context, id uint) error
	CreateTemplate(ctx context.Context, file []byte, template entities.Template) (entities.Template, error)
	DeleteTemplate(ctx context.Context, id uint) error
}

func NewMemeService(uowFactory repositories.UoWFactory) Meme {
	return &memeService{uowFactory: uowFactory}
}

type memeService struct {
	uowFactory repositories.UoWFactory
}

func (m *memeService) GetAll(ctx context.Context) ([]entities.Meme, error) {
	uow := m.uowFactory.Create()
	memeRepo := uow.GetMemeRepository()

	memes, err := memeRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return memes, nil
}

func (m *memeService) Get(ctx context.Context, id uint) (entities.Meme, error) {
	uow := m.uowFactory.Create()
	memeRepo := uow.GetMemeRepository()

	meme, err := memeRepo.Get(id)
	if err != nil {
		return meme, err
	}
	return meme, nil
}

func (m *memeService) GetByAuthor(ctx context.Context, userId uint) ([]entities.Meme, error) {
	uow := m.uowFactory.Create()
	memeRepo := uow.GetMemeRepository()

	meme, err := memeRepo.GetByAuthor(userId)
	if err != nil {
		return meme, err
	}
	return meme, nil
}

func (m *memeService) Create(ctx context.Context, file []byte, meme entities.Meme) (entities.Meme, error) {
	uow := m.uowFactory.Create()

	err := uow.BeginTransaction()
	if err != nil {
		return entities.Meme{}, err
	}

	memeRepo := uow.GetMemeRepository()
	id, err := memeRepo.Create(meme)
	if err != nil {
		err := tryRollback(uow, err)
		return entities.Meme{}, err
	}

	fileRepo := uow.GetFileRepository()
	err = fileRepo.Save(file, meme.FilePath)
	if err != nil {
		err := tryRollback(uow, err)
		return entities.Meme{}, err
	}

	err = uow.CommitTransaction()
	if err != nil {
		err := tryRollback(uow, err)
		return entities.Meme{}, err
	}

	memeRepo = uow.GetMemeRepository()
	memeResult, err := memeRepo.Get(id)
	if err != nil {
		return entities.Meme{}, err
	}
	return memeResult, nil
}

func tryRollback(uow repositories.UnitOfWork, err error) error {
	rollbackErr := uow.RollbackTransaction()
	if rollbackErr != nil {
		return rollbackErr
	}
	return err
}

func (m *memeService) AddComment(ctx context.Context, memeID uint, comment entities.Comment) (entities.Meme, error) {
	panic("implement me")
}

func (m *memeService) DeleteComment(ctx context.Context, memeID uint, commentId uint) (entities.Meme, error) {
	//TODO: check rights
	panic("implement me")
}

func (m *memeService) Delete(ctx context.Context, id uint) error {
	//TODO: check rights
	uow := m.uowFactory.Create()
	memeRepo := uow.GetMemeRepository()

	err := memeRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (m *memeService) CreateTemplate(ctx context.Context, file []byte, template entities.Template) (entities.Template, error) {
	//TODO: check rights
	panic("implement me")
}

func (m *memeService) DeleteTemplate(ctx context.Context, id uint) error {
	//TODO: check rights
	panic("implement me")
}
