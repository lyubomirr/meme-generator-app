package persistence

import (
	"errors"
	"fmt"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	customErr "github.com/lyubomirr/meme-generator-app/core/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	usernameMaxLength = 25
	minPasswordLength = 8
)

type dbUser struct {
	ID         uint
	Username   string `gorm:"type:varchar(25);uniqueIndex"`
	Password   string
	RoleID     uint
	Role       dbRole
	PictureURL string
	Memes      []dbMeme    `gorm:"foreignKey:AuthorID"`
	Comment    []dbComment `gorm:"foreignKey:AuthorID"`
}

func (dbUser) TableName() string {
	return "users"
}

func (u dbUser) toEntity() entities.User {
	return entities.User{
		ID:         u.ID,
		Username:   u.Username,
		Password:   u.Password,
		Role:       u.Role.toEntity(),
		PictureURL: u.PictureURL,
	}
}

func newUser(entity entities.User) dbUser {
	return dbUser{
		ID:         entity.ID,
		Username:   entity.Username,
		Password:   entity.Password,
		RoleID:     entity.Role.ID,
		PictureURL: entity.PictureURL,
	}
}

type mySqlUserRepository struct {
	db *gorm.DB
}

func (m *mySqlUserRepository) Get(id uint) (entities.User, error) {
	var user dbUser
	result := withReadTimeout(m.db).
		Preload(clause.Associations).
		First(&user, id)

	if result.Error != nil {
		return entities.User{}, result.Error
	}
	return user.toEntity(), nil
}

func (m *mySqlUserRepository) GetByUsername(username string) (entities.User, error) {
	var user dbUser
	result := withReadTimeout(m.db).
		Preload(clause.Associations).
		Where(&dbUser{Username: username}).
		First(&user)

	if result.Error != nil {
		return entities.User{}, result.Error
	}
	return user.toEntity(), nil
}

func (m *mySqlUserRepository) Create(user entities.User) (uint, error) {
	err := checkUserConstraints(user)
	if err != nil {
		return 0, err
	}

	_, err = m.GetByUsername(user.Username)
	if err == nil {
		return 0, customErr.ExistingResourceError{
			Err: errors.New(fmt.Sprintf("User with name %v already exists", user.Username)),
		}
	}

	model := newUser(user)
	hashed, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("hashing password failed: %w", err)
	}

	model.Password = string(hashed)

	result := m.db.Create(&model)
	if result.Error != nil {
		return 0, result.Error
	}
	return model.ID, nil
}

func checkUserConstraints(user entities.User) error {
	if len(user.Username) > usernameMaxLength {
		return customErr.NewValidationError(
			fmt.Errorf("username length cannot contain more than %v symbols", usernameMaxLength))
	}

	if len(user.Password) < minPasswordLength {
		return customErr.NewValidationError(
			fmt.Errorf(fmt.Sprintf("password should contain at least %v symbols", minPasswordLength)))
	}

	return nil
}

func (m *mySqlUserRepository) Update(user entities.User) (entities.User, error) {
	var dbUser dbUser
	result := m.db.First(&dbUser, user.ID)
	if result.Error != nil {
		return entities.User{}, result.Error
	}

	dbUser.RoleID = user.Role.ID
	dbUser.PictureURL = user.PictureURL

	if dbUser.Password != user.Password {
		//A new password is set - validate length and hash
		if len(user.Password) < minPasswordLength {
			return entities.User{},
				customErr.NewValidationError(
					fmt.Errorf("password should contain at least %v symbols", minPasswordLength))
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return entities.User{}, fmt.Errorf("hashing password failed: %w", err)
		}
		dbUser.Password = string(hashed)
	}

	result = m.db.Save(dbUser)
	if result.Error != nil {
		return entities.User{}, result.Error
	}

	return user, nil
}

func (m *mySqlUserRepository) Delete(id uint) error {
	_, err := m.Get(id)
	if err != nil {
		return err
	}

	result := m.db.Delete(&dbUser{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
