package persistence

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	"gorm.io/gorm"
	"time"
)

type dbMeme struct {
	ID         uint
	AuthorID   uint
	Author     dbUser
	Title      string `gorm:"type:varchar(50)"`
	FilePath   string
	MimeType   string      `gorm:"type:varchar(50)"`
	Comments   []dbComment `gorm:"foreignKey:MemeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TemplateID sql.NullInt64
	Template   dbTemplate `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time
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
		AuthorID: m.AuthorID,
		Author:   m.Author.toEntity(),
		Title:    m.Title,
		FilePath: m.FilePath,
		MimeType: m.MimeType,
		Comments: comments,
		TemplateID: uint(m.TemplateID.Int64),
		Template: m.Template.toEntity(),
		CreatedAt: m.CreatedAt,
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
		AuthorID: meme.AuthorID,
		Title:    meme.Title,
		FilePath: meme.FilePath,
		MimeType: meme.MimeType,
		Comments: commentsToDbModels(meme.Comments),
		TemplateID: sql.NullInt64{
			Int64: int64(meme.TemplateID),
			Valid: meme.TemplateID != 0,
		},
		CreatedAt: meme.CreatedAt,
	}
}

type mySqlMemeRepository struct {
	db       *gorm.DB
	validate *validator.Validate
}

func (m mySqlMemeRepository) GetAll() ([]entities.Meme, error) {
	var memes []dbMeme
	result := preloadMemes(m.db).Order("created_at desc").Find(&memes)

	if result.Error != nil {
		return nil, result.Error
	}

	return memesToEntities(memes), nil
}

func preloadMemes(db *gorm.DB) *gorm.DB {
	return withReadTimeout(db).
		Preload("Author").
		Preload("Template").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("comments.created_at desc")
		}).
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
	err := m.validate.Struct(meme)
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

func (m *mySqlMemeRepository) Update(meme entities.Meme) (entities.Meme, error) {
	err := m.validate.Struct(meme)
	if err != nil {
		return entities.Meme{}, err
	}

	var dbMeme dbMeme
	result := preloadMemes(m.db).First(&dbMeme, meme.ID)
	if result.Error != nil {
		return entities.Meme{}, result.Error
	}

	dbMeme = newMeme(meme)
	m.db.Model(&dbMeme).Association("Comments").Replace(dbMeme.Comments)

	result = m.db.Save(&dbMeme)
	if result.Error != nil {
		return entities.Meme{}, result.Error
	}

	updated, err := m.Get(meme.ID)
	if err != nil {
		return entities.Meme{}, err
	}
	return updated, nil
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
