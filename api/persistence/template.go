package persistence

import (
	"errors"
	"fmt"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	customErr "github.com/lyubomirr/meme-generator-app/core/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	memeTemplateNameMaxLength = 50
)

type dbTemplateTextPosition struct {
	TopOffset uint `gorm:"primaryKey"`
	LeftOffset uint `gorm:"primaryKey"`
	TemplateID uint `gorm:"primaryKey"`
}

func (dbTemplateTextPosition) TableName() string {
	return "template_text_positions"
}

type dbTemplate struct {
	ID       uint
	Name    string `gorm:"type:varchar(50)"`
	FilePath string
	TextPositions []dbTemplateTextPosition `gorm:"foreignKey:TemplateID"`
}

func (dbTemplate) TableName() string {
	return "templates"
}

func textPositionsToEntities(positions []dbTemplateTextPosition) []entities.TemplateTextPosition {
	e := make([]entities.TemplateTextPosition, 0, len(positions))
	for _, p := range positions {
		e = append(e, entities.TemplateTextPosition{
			TopOffset:  p.TopOffset,
			LeftOffset: p.LeftOffset,
		})
	}
	return e
}

func (m dbTemplate) toEntity() entities.Template {
	return entities.Template{
		ID:       m.ID,
		Name:     m.Name,
		FilePath: m.FilePath,
		TextPositions: textPositionsToEntities(m.TextPositions),
	}
}

func newTemplate(template entities.Template) dbTemplate {
	return dbTemplate{
		ID:       template.ID,
		Name:     template.Name,
		FilePath: template.FilePath,
	}
}

type mySqlTemplateRepository struct {
	db *gorm.DB
}

func (m mySqlTemplateRepository) GetAll() ([]entities.Template, error) {
	var templates []dbTemplate
	result := withReadTimeout(m.db).Preload(clause.Associations).Find(&templates)

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
	result := withReadTimeout(m.db).Preload(clause.Associations).First(&template, id)

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
	if template.Name == "" {
		return customErr.NewValidationError(errors.New("template name is empty"))
	}
	if template.FilePath == "" {
		return customErr.NewValidationError(errors.New("filepath is empty"))
	}
	if len(template.TextPositions) == 0 {
		return customErr.NewValidationError(errors.New("you must specify at least one position for inserting text"))
	}
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
