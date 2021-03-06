package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	customErr "github.com/lyubomirr/meme-generator-app/core/errors"
	"github.com/lyubomirr/meme-generator-app/core/repositories"
	"path"
)

type Template interface {
	GetAll(ctx context.Context) ([]entities.Template, error)
	Get(ctx context.Context, id uint) (entities.Template, error)
	Create(ctx context.Context, file []byte, template entities.Template) (entities.Template, error)
	Delete(ctx context.Context, id uint) error
}

func NewTemplateService(uowFactory repositories.UoWFactory) Template {
	return &templateService{uowFactory: uowFactory}
}

type templateService struct {
	uowFactory repositories.UoWFactory
}

func (t *templateService) Get(ctx context.Context, id uint) (entities.Template, error) {
	uow := t.uowFactory.Create()
	templates := uow.GetTemplateRepository()
	template, err := templates.Get(id)
	if err != nil {
		return entities.Template{}, err
	}
	return template, nil
}

func (t *templateService) GetAll(ctx context.Context) ([]entities.Template, error) {
	uow := t.uowFactory.Create()
	templates := uow.GetTemplateRepository()
	temp, err := templates.GetAll()
	if err != nil {
		return nil, err
	}
	return temp, nil
}

func (t *templateService) Create(
	ctx context.Context, file []byte, template entities.Template) (result entities.Template, err error) {

	if !IsAdministrator(ctx) {
		return entities.Template{}, customErr.NewRightsError(errors.New("user is not administrator"))
	}

	uow := t.uowFactory.Create()
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

	mimeType, err := getMimeType(file)
	if err != nil {
		return entities.Template{}, err
	}
	fileName := fmt.Sprintf("%v%v", uuid.NewString(), getFileExtension(mimeType))

	template.MimeType = mimeType
	template.FilePath = path.Join(templateFilesPath, fileName)

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

func (t *templateService) Delete(ctx context.Context, id uint) (err error) {
	if !IsAdministrator(ctx) {
		return customErr.NewRightsError(errors.New("user is not administrator"))
	}

	uow := t.uowFactory.Create()
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