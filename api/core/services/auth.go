package services

import (
	"errors"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	customErr "github.com/lyubomirr/meme-generator-app/core/errors"
	"github.com/lyubomirr/meme-generator-app/core/repositories"
	"golang.org/x/crypto/bcrypt"
)

type LoginModel struct {
	Username string
	Password string
}

type NewRegistrationModel struct {
	LoginModel
	ConfirmPassword string
}

type Authentication interface {
	ValidateCredentials(model LoginModel) error
	Register(model NewRegistrationModel) (entities.User, error)
}

func NewAuthService(userRepo repositories.User) Authentication {
	return &authService{userRepository: userRepo}
}

type authService struct {
	userRepository repositories.User
}

func (a authService) ValidateCredentials(model LoginModel) error {
	user, err := a.userRepository.GetByUsername(model.Username)
	if err != nil {
		return customErr.AuthError{
			Err: err,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(model.Password))
	if err != nil {
		return customErr.AuthError{
			Err: errors.New("incorrect password"),
		}
	}
	return nil
}

func (a authService) Register(model NewRegistrationModel) (entities.User, error) {
	panic("implement me")
}

