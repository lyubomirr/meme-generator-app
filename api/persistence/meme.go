package persistence

import (
	"fmt"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	customErr "github.com/lyubomirr/meme-generator-app/core/errors"
	"github.com/lyubomirr/meme-generator-app/core/repositories"
	"gorm.io/gorm"
)

const (
	memeNameMaxLength = 50
)

type dbMeme struct {
	ID       uint
	AuthorID uint
	Author   dbUser
	Title    string `gorm:"type:varchar(50)"`
	FilePath string
	Comments []dbComment `gorm:"foreignKey:MemeID"`
}

func (dbMeme) TableName() string {
	return "memes"
}

func (m dbMeme) toEntity() entities.Meme {
	comments := make([]entities.Comment, 0, len(m.Comments))
	for _, c := range m.Comments {
		comments = append(comments, c.toEntity())
	}

	return entities.Meme{
		ID:       m.ID,
		Author:   m.Author.toEntity(),
		Title:    m.Title,
		FilePath: m.FilePath,
		Comments: comments,
	}
}

func memesToEntities(dbMemes []dbMeme) []entities.Meme {
	memes := make([]entities.Meme, 0, len(dbMemes))
	for _, m := range dbMemes {
		memes = append(memes, m.toEntity())
	}
	return memes
}

func newMeme(meme entities.Meme) dbMeme {
	return dbMeme{
		ID:       meme.ID,
		AuthorID: meme.Author.ID,
		Title:    meme.Title,
		FilePath: meme.FilePath,
		Comments: commentsToDbModels(meme.Comments),
	}
}

func NewMemeRepository() repositories.Meme {
	return &mySqlMemeRepository{db: getDB()}
}

type mySqlMemeRepository struct {
	db *gorm.DB
}

func (m mySqlMemeRepository) GetAll() ([]entities.Meme, error) {
	var memes []dbMeme
	result := preloadMemes(m.db).Find(&memes)

	if result.Error != nil {
		return nil, result.Error
	}

	return memesToEntities(memes), nil
}

func preloadMemes(db *gorm.DB) *gorm.DB {
	return withReadTimeout(db).
		Preload("Author").
		Preload("Comments").
		Preload("Comments.Author")
}

func (m *mySqlMemeRepository) Get(id uint) (entities.Meme, error) {
	var meme dbMeme
	result := preloadMemes(m.db).First(&meme, id)

	if result.Error != nil {
		return entities.Meme{}, result.Error
	}
	return meme.toEntity(), nil
}

func (m *mySqlMemeRepository) GetByAuthor(userId uint) ([]entities.Meme, error) {
	var memes []dbMeme
	result := preloadMemes(m.db).
		Where(&dbMeme{AuthorID: userId}).
		Find(memes)

	if result.Error != nil {
		return nil, result.Error
	}
	return memesToEntities(memes), nil
}

func (m *mySqlMemeRepository) Create(meme entities.Meme) (uint, error) {
	err := checkMemeConstraints(meme)
	if err != nil {
		return 0, err
	}

	model := newMeme(meme)
	result := m.db.Create(&model)
	if result.Error != nil {
		return 0, result.Error
	}
	return model.ID, nil
}

func checkMemeConstraints(meme entities.Meme) error {
	if len(meme.Title) > memeNameMaxLength {
		return customErr.NewValidationError(
			fmt.Errorf("meme title cannot contain more than %v symbols", usernameMaxLength))
	}
	return nil
}

func (m *mySqlMemeRepository) Update(meme entities.Meme) (entities.Meme, error) {
	err := checkMemeConstraints(meme)
	if err != nil {
		return entities.Meme{}, err
	}

	var dbMeme dbMeme
	result := preloadMemes(m.db).First(&dbMeme, meme.ID)
	if result.Error != nil {
		return entities.Meme{}, result.Error
	}

	dbMeme = newMeme(meme)
	result = m.db.Save(dbMeme)
	if result.Error != nil {
		return entities.Meme{}, result.Error
	}
	return dbMeme.toEntity(), nil
}

func (m mySqlMemeRepository) Delete(id uint) error {
	_, err := m.Get(id)
	if err != nil {
		return err
	}

	result := m.db.Delete(&dbMeme{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}