package persistence

import (
	"fmt"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	customErr "github.com/lyubomirr/meme-generator-app/core/errors"
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
	var templates []dbTemplate
	result := withReadTimeout(m.db).Find(&templates)

	if result.Error != nil {
		return nil, result.Error
	}

	var entities = make([]entities.Template, 0, len(templates))
	for _, t := range templates {
		entities = append(entities, t.toEntity())
	}

	return entities, nil
}

func (m mySqlTemplateRepository) Get(id uint) (entities.Template, error) {
	var template dbTemplate
	result := withReadTimeout(m.db).First(&template, id)

	if result.Error != nil {
		return entities.Template{}, result.Error
	}
	return template.toEntity(), nil
}

func (m mySqlTemplateRepository) Create(template entities.Template) (uint, error) {
	err := checkTemplateConstraints(template)
	if err != nil {
		return 0, err
	}

	model := newTemplate(template)
	result := m.db.Create(&model)
	if result.Error != nil {
		return 0, result.Error
	}
	return model.ID, nil
}

func checkTemplateConstraints(template entities.Template) error {
	if len(template.Name) > memeTemplateNameMaxLength {
		return customErr.NewValidationError(
			fmt.Errorf("template name cannot contain more than %v symbols", memeTemplateNameMaxLength))
	}
	return nil
}

func (m mySqlTemplateRepository) Update(template entities.Template) (entities.Template, error) {
	err := checkTemplateConstraints(template)
	if err != nil {
		return entities.Template{}, err
	}

	var dbTemplate dbTemplate
	result := withReadTimeout(m.db).First(&dbTemplate, template.ID)
	if result.Error != nil {
		return entities.Template{}, result.Error
	}

	dbTemplate = newTemplate(template)
	result = m.db.Save(dbTemplate)
	if result.Error != nil {
		return entities.Template{}, result.Error
	}
	return dbTemplate.toEntity(), nil
}

func (m mySqlTemplateRepository) Delete(id uint) error {
	_, err := m.Get(id)
	if err != nil {
		return err
	}

	result := m.db.Delete(&dbTemplate{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
