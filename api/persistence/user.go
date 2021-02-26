package persistence

import (
	"fmt"
	"github.com/lyubomirr/meme-generator/core/entities"
	"github.com/lyubomirr/meme-generator/core/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

const (
	usernameMaxLength = 25
	minPasswordLength = 8
)

type dbUser struct {
	gorm.Model
	Username   string `gorm:"type:varchar(25);uniqueIndex"`
	Password   string
	RoleID     uint
	Role       dbRole
	PictureURL string
}

func (dbUser) TableName() string {
	return "users"
}

func (u dbUser) toEntity() entities.User {
	return entities.User{
		Username: u.Username,
		Password: u.Password,
		Role: u.Role.toEntity(),
		PictureURL: u.PictureURL,
	}
}

func newUser(entity entities.User) dbUser {
	return dbUser{
		Model:      gorm.Model{
			ID:        entity.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
		Username:   entity.Username,
		Password:   entity.Password,
		RoleID:     entity.Role.ID,
		PictureURL: entity.PictureURL,
	}
}

func NewUserRepository() repositories.User {
	return &mySqlUserRepository{db: getDB()}
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

func (m *mySqlUserRepository) Create(user entities.User) (entities.User, error) {
	err := checkUserConstraints(user)
	if err != nil {
		return entities.User{}, err
	}

	u, err := m.GetByUsername(user.Username)
	if err == nil && u != (entities.User{}) {
		return entities.User{}, fmt.Errorf("the user already exists")
	}

	model := newUser(user)
	hashed, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.User{}, fmt.Errorf("hashing password failed: %w", err)
	}

	model.Password = string(hashed)

	result := m.db.Create(model)
	if result.Error != nil {
		return entities.User{}, result.Error
	}
	return user, nil
}

func checkUserConstraints(user entities.User) error {
	if len(user.Username) > usernameMaxLength {
		return fmt.Errorf("username length cannot contain more than %v symbols", usernameMaxLength)
	}

	if len(user.Password) < minPasswordLength {
		return fmt.Errorf(fmt.Sprintf("password should contain at least %v symbols", minPasswordLength))
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
	dbUser.UpdatedAt = time.Now()

	if dbUser.Password != user.Password {
		//A new password is set - validate length and hash
		if len(user.Password) < minPasswordLength {
			return entities.User{},
			fmt.Errorf(fmt.Sprintf("password should contain at least %v symbols", minPasswordLength))
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


