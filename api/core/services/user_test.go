package services

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	customErr "github.com/lyubomirr/meme-generator-app/core/errors"
	"github.com/lyubomirr/meme-generator-app/core/mocks"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"testing"
)

func TestUserService_ValidateCredentials(t *testing.T) {
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

	sut := userService{userRepository: m}
	_, err = sut.ValidateCredentials(user.Username, plainTextPassword)

	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
}

func TestUserService_ValidateCredentials_ShouldReturnErrIfNoUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidUser := "pesho"

	m := mocks.NewMockUserRepository(ctrl)
	m.EXPECT().GetByUsername(gomock.Eq(invalidUser)).Return(entities.User{}, errors.New("no user"))

	sut := userService{userRepository: m}
	_, err := sut.ValidateCredentials(invalidUser, "alabala")

	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestUserService_ValidateCredentials_ShouldReturnErrIfWrongPassword(t *testing.T) {
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

	sut := userService{userRepository: m}
	_, err := sut.ValidateCredentials(user.Username, plainTextPassword)

	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := entities.User{
		ID:         1,
		Username:   "admin",
		Password:   "admin",
		Role:       entities.Role{
			ID: 4,
			Name: "alabala",
		},
		PictureURL: "",
	}

	m := mocks.NewMockUserRepository(ctrl)
	m.EXPECT().GetByUsername(gomock.Eq(user.Username)).Return(entities.User{}, gorm.ErrRecordNotFound)
	m.EXPECT().Create(gomock.Eq(user)).Return(user.ID, nil)
	m.EXPECT().Get(gomock.Eq(user.ID)).Return(user, nil)

	sut := userService{userRepository: m}
	_, err := sut.Create(user)

	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
}

func TestUserService_Create_ShouldReturnErrIfUserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := entities.User{
		ID:         1,
		Username:   "admin",
		Password:   "admin",
		Role:       entities.Role{
			ID: 4,
			Name: "alabala",
		},
		PictureURL: "",
	}

	m := mocks.NewMockUserRepository(ctrl)
	m.EXPECT().GetByUsername(gomock.Eq(user.Username)).
		Return(entities.User{ ID: 4, Username: "admin"}, nil)

	sut := userService{userRepository: m}
	_, err := sut.Create(user)

	if err == nil || !errors.As(err, &customErr.ExistingResourceError{}) {
		t.Errorf("expected error but got nil or different one")
	}
}