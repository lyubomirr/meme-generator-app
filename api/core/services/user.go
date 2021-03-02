package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	customErr "github.com/lyubomirr/meme-generator-app/core/errors"
	"github.com/lyubomirr/meme-generator-app/core/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User interface {
	ValidateCredentials(ctx context.Context, username string, password string) (entities.User, error)
	Create(ctx context.Context, user entities.User) (entities.User, error)
}

func NewUserService(uowFactory repositories.UoWFactory) User {
	return &userService{uowFactory: uowFactory}
}

type userService struct {
	uowFactory repositories.UoWFactory
}

func (a *userService) ValidateCredentials(
	ctx context.Context, username string, password string) (entities.User, error) {
	uow := a.uowFactory.Create()
	repo := uow.GetUserRepository()

	user, err := repo.GetByUsername(username)
	if err != nil {
		return entities.User{}, customErr.NewAuthError(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return entities.User{}, customErr.NewAuthError(errors.New("incorrect password"))
	}
	return user, nil
}

func (a *userService) Create(ctx context.Context, user entities.User) (entities.User, error) {
	uow := a.uowFactory.Create()
	repo := uow.GetUserRepository()

	u, err := repo.GetByUsername(user.Username)
	if err == nil {
		return entities.User{}, customErr.NewExistingResourceError(
			errors.New(fmt.Sprintf("User with name %v already exists", user.Username)))
	}

	//TODO: ADD CHECK FOR ADMIN
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.User{}, err
	}

	id, err := repo.Create(user)
	if err != nil {
		return entities.User{}, err
	}

	u, err = repo.Get(id)
	if err != nil {
		return entities.User{}, err
	}
	return u, nil
}

