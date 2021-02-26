package services

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	"github.com/lyubomirr/meme-generator-app/core/mocks"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestAuthService_ValidateCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	plainTextPassword := "admin123$%$"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Error("couldnt hash password")
	}

	user := entities.User{
		ID:         1,
		Username:   "admin",
		Password:   string(hashedPassword),
		Role:       entities.Role{},
		PictureURL: "",
	}

	m := mocks.NewMockUserRepository(ctrl)
	m.EXPECT().GetByUsername(gomock.Eq(user.Username)).Return(user, nil)

	sut := authService{userRepository: m}
	err = sut.ValidateCredentials(LoginModel{
		Username: user.Username,
		Password: plainTextPassword,
	})

	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
}

func TestAuthService_ValidateCredentials_ShouldFailIfNoUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidUser := "pesho"

	m := mocks.NewMockUserRepository(ctrl)
	m.EXPECT().GetByUsername(gomock.Eq(invalidUser)).Return(entities.User{}, errors.New("no user"))

	sut := authService{userRepository: m}
	err := sut.ValidateCredentials(LoginModel{
		Username: invalidUser,
		Password: "alabala",
	})

	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestAuthService_ValidateCredentials_ShouldFailIfWrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	plainTextPassword := "admin123$%$"
	user := entities.User{
		ID:         1,
		Username:   "admin",
		Password:   "differenthash",
		Role:       entities.Role{},
		PictureURL: "",
	}

	m := mocks.NewMockUserRepository(ctrl)
	m.EXPECT().GetByUsername(gomock.Eq(user.Username)).Return(user, nil)

	sut := authService{userRepository: m}
	err := sut.ValidateCredentials(LoginModel{
		Username: user.Username,
		Password: plainTextPassword,
	})

	if err == nil {
		t.Errorf("expected error but got nil")
	}
}