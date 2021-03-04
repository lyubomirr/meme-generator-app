package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	customErr "github.com/lyubomirr/meme-generator-app/core/errors"
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

func (m *memeService) Create(ctx context.Context, file []byte, meme entities.Meme) (result entities.Meme, err error) {
	uow := m.uowFactory.Create()
	defer func() {
		if r := recover(); r != nil {
			err = uow.RollbackTransaction()
			if err != nil {
				err = fmt.Errorf("failed to rollback on panic: %v", r)
			} else  {
				err = fmt.Errorf("panic: %v", r)
			}
		}
	}()

	err = uow.BeginTransaction()
	if err != nil {
		return
	}

	memeRepo := uow.GetMemeRepository()
	id, err := memeRepo.Create(meme)
	if err != nil {
		err = tryRollback(uow, err)
		return
	}

	fileRepo := uow.GetFileRepository()
	err = fileRepo.Save(file, meme.FilePath)
	if err != nil {
		err = tryRollback(uow, err)
		return
	}

	err = uow.CommitTransaction()
	if err != nil {
		err = tryRollback(uow, err)
		return
	}

	memeRepo = uow.GetMemeRepository()
	result, err = memeRepo.Get(id)
	return
}

func tryRollback(uow repositories.UnitOfWork, err error) error {
	rollbackErr := uow.RollbackTransaction()
	if rollbackErr != nil {
		return rollbackErr
	}
	return err
}

func (m *memeService) AddComment(ctx context.Context, memeID uint, comment entities.Comment) (entities.Meme, error) {
	uow := m.uowFactory.Create()
	memes := uow.GetMemeRepository()

	meme, err := memes.Get(memeID)
	if err != nil {
		return entities.Meme{}, err
	}

	meme.Comments = append(meme.Comments, comment)
	meme, err = memes.Update(meme)
	if err != nil {
		return entities.Meme{}, err
	}
	return meme, nil
}

func (m *memeService) DeleteComment(ctx context.Context, memeID uint, commentId uint) (entities.Meme, error) {
	//TODO: check rights
	uow := m.uowFactory.Create()
	memes := uow.GetMemeRepository()

	meme, err := memes.Get(memeID)
	if err != nil {
		return entities.Meme{}, err
	}

	commentIdx := -1
	for idx, c := range meme.Comments {
		if c.ID == commentId {
			commentIdx = idx
			break
		}
	}

	if commentIdx == -1 {
		return entities.Meme{}, customErr.NewValidationError(errors.New("no such comment"))
	}

	meme.Comments = append(meme.Comments[:commentIdx], meme.Comments[commentIdx+1:]...)
	meme, err = memes.Update(meme)
	if err != nil {
		return entities.Meme{}, err
	}
	return meme, nil
}

func (m *memeService) Delete(ctx context.Context, id uint) (err error) {
	//TODO: check rights
	uow := m.uowFactory.Create()
	defer func() {
		if r := recover(); r != nil {
			err = uow.RollbackTransaction()
			if err != nil {
				err = fmt.Errorf("failed to rollback on panic: %v", r)
			} else  {
				err = fmt.Errorf("panic: %v", r)
			}
		}
	}()

	err = uow.BeginTransaction()
	if err != nil {
		return
	}

	memes := uow.GetMemeRepository()
	meme, err := memes.Get(id)
	if err != nil {
		return tryRollback(uow, err)
	}

	err = memes.Delete(id)
	if err != nil {
		return tryRollback(uow, err)
	}

	fileRepo := uow.GetFileRepository()
	err = fileRepo.Delete(meme.FilePath)
	if err != nil {
		return tryRollback(uow, err)
	}

	err = uow.CommitTransaction()
	if err != nil {
		return tryRollback(uow, err)
	}
	return
}

func (m *memeService) CreateTemplate(
	ctx context.Context, file []byte, template entities.Template) (result entities.Template, err error) {

	if !isAdministrator(ctx) {
		return entities.Template{}, customErr.NewRightsError(errors.New("user is not administrator"))
	}

	uow := m.uowFactory.Create()
	defer func() {
		if r := recover(); r != nil {
			err = uow.RollbackTransaction()
			if err != nil {
				err = fmt.Errorf("failed to rollback on panic: %v", r)
			} else  {
				err = fmt.Errorf("panic: %v", r)
			}
		}
	}()

	err = uow.BeginTransaction()
	if err != nil {
		return
	}

	templates := uow.GetTemplateRepository()
	id, err := templates.Create(template)
	if err != nil {
		err = tryRollback(uow, err)
		return
	}

	fileRepo := uow.GetFileRepository()
	err = fileRepo.Save(file, template.FilePath)
	if err != nil {
		err = tryRollback(uow, err)
		return
	}

	err = uow.CommitTransaction()
	if err != nil {
		err = tryRollback(uow, err)
		return
	}

	templates = uow.GetTemplateRepository()
	result, err = templates.Get(id)
	return
}

func (m *memeService) DeleteTemplate(ctx context.Context, id uint) (err error) {
	//TODO: check rights
	uow := m.uowFactory.Create()
	defer func() {
		if r := recover(); r != nil {
			err = uow.RollbackTransaction()
			if err != nil {
				err = fmt.Errorf("failed to rollback on panic: %v", r)
			} else  {
				err = fmt.Errorf("panic: %v", r)
			}
		}
	}()

	err = uow.BeginTransaction()
	if err != nil {
		return err
	}

	templates := uow.GetTemplateRepository()
	template, err := templates.Get(id)
	if err != nil {
		return tryRollback(uow, err)
	}

	err = templates.Delete(id)
	if err != nil {
		return tryRollback(uow, err)
	}
	//return tryRollback(uow, err)

	fileRepo := uow.GetFileRepository()
	err = fileRepo.Delete(template.FilePath)
	if err != nil {
		return tryRollback(uow, err)
	}

	err = uow.CommitTransaction()
	if err != nil {
		return tryRollback(uow, err)
	}
	return
}
