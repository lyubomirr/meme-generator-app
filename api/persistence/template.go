package persistence

import (
	"github.com/lyubomirr/meme-generator-app/core/entities"
	"github.com/lyubomirr/meme-generator-app/core/repositories"
	"gorm.io/gorm"
)

const (
	memeTemplateNameMaxLength = 50
)

type dbTemplate struct {
	ID       uint
	Name    string `gorm:"type:varchar(50)"`
	FilePath string
}

func (dbTemplate) TableName() string {
	return "templates"
}

func (m dbTemplate) toEntity() entities.Template {
	return entities.Template{
		ID:       m.ID,
		Name:     m.Name,
		FilePath: m.FilePath,
	}
}

func newTemplate(template entities.Template) dbTemplate {
	return dbTemplate{
		ID:       template.ID,
		Name:     template.Name,
		FilePath: template.FilePath,
	}
}

func NewTemplateRepository() repositories.Template {
	return &mySqlTemplateRepository{db: getDB()}
}

type mySqlTemplateRepository struct {
	db *gorm.DB
}

func (m mySqlTemplateRepository) GetAll() ([]entities.Template, error) {
	panic("implement me")
}

func (m mySqlTemplateRepository) Get(id uint) (entities.Template, error) {
	panic("implement me")
}

func (m mySqlTemplateRepository) Create(meme entities.Template) (uint, error) {
	panic("implement me")
}

func (m mySqlTemplateRepository) Update(meme entities.Template) (entities.Template, error) {
	panic("implement me")
}

func (m mySqlTemplateRepository) Delete(id uint) error {
	panic("implement me")
}
