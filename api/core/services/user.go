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

func NewUserService(userRepo repositories.User) User {
	return &userService{userRepository: userRepo}
}

type userService struct {
	userRepository repositories.User
}

func (a *userService) ValidateCredentials(
	ctx context.Context, username string, password string) (entities.User, error) {
	user, err := a.userRepository.GetByUsername(username)
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
	u, err := a.userRepository.GetByUsername(user.Username)
	if err == nil {
		return entities.User{}, customErr.NewExistingResourceError(
			errors.New(fmt.Sprintf("User with name %v already exists", user.Username)))
	}

	//TODO: ADD CHECK FOR ADMIN
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.User{}, err
	}

	id, err := a.userRepository.Create(user)
	if err != nil {
		return entities.User{}, err
	}

	u, err = a.userRepository.Get(id)
	if err != nil {
		return entities.User{}, err
	}
	return u, nil
}

