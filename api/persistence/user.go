package persistence

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	customErr "github.com/lyubomirr/meme-generator-app/core/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type dbUser struct {
	ID         uint
	Username   string `gorm:"type:varchar(25);uniqueIndex"`
	Password   string
	RoleID     uint
	Role       dbRole
	PictureURL string
	Memes      []dbMeme    `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Comments   []dbComment `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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
	}
}

func newUser(entity entities.User) dbUser {
	return dbUser{
		ID:         entity.ID,
		Username:   entity.Username,
		Password:   entity.Password,
		RoleID:     entity.Role.ID,
	}
}

type mySqlUserRepository struct {
	db *gorm.DB
	validate *validator.Validate
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
	err := m.validate.Struct(user)
	if err != nil {
		return 0, err
	}

	_, err = m.GetByUsername(user.Username)
	if err == nil {
		return 0, customErr.ValidationError{
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

func (m *mySqlUserRepository) Update(user entities.User) (entities.User, error) {
	var dbUser dbUser
	result := m.db.First(&dbUser, user.ID)
	if result.Error != nil {
		return entities.User{}, result.Error
	}

	err := m.validate.StructPartial(user, "Role.ID")
	if err != nil {
		return entities.User{}, err
	}

	dbUser.RoleID = user.Role.ID

	if user.Password != "" && dbUser.Password != user.Password {
		//A new password is set - validate length and hash
		err = m.validate.StructPartial(user, "Password")
		if err != nil {
			return entities.User{}, err
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return entities.User{}, fmt.Errorf("hashing password failed: %w", err)
		}
		dbUser.Password = string(hashed)
	}

	result = m.db.Save(&dbUser)
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
