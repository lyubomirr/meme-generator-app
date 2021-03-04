package persistence

import (
	"github.com/go-playground/validator/v10"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	MimeType string `gorm:"type:varchar(50)"`
	TextPositions []dbTemplateTextPosition `gorm:"foreignKey:TemplateID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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

func newTextPositions(positions []entities.TemplateTextPosition) []dbTemplateTextPosition {
	dbModels := make([]dbTemplateTextPosition, 0, len(positions))
	for _, p := range positions {
		dbModels = append(dbModels, dbTemplateTextPosition{
			LeftOffset: p.LeftOffset,
			TopOffset: p.TopOffset,
		})
	}
	return dbModels
}

func (m dbTemplate) toEntity() entities.Template {
	return entities.Template{
		ID:       m.ID,
		Name:     m.Name,
		FilePath: m.FilePath,
		MimeType: m.MimeType,
		TextPositions: textPositionsToEntities(m.TextPositions),
	}
}

func newTemplate(template entities.Template) dbTemplate {
	return dbTemplate{
		ID:       template.ID,
		Name:     template.Name,
		FilePath: template.FilePath,
		MimeType: template.MimeType,
		TextPositions: newTextPositions(template.TextPositions),
	}
}

type mySqlTemplateRepository struct {
	db *gorm.DB
	validate *validator.Validate
}

func (m mySqlTemplateRepository) GetAll() ([]entities.Template, error) {
	var templates []dbTemplate
	result := withReadTimeout(m.db).Preload(clause.Associations).Find(&templates)

	if result.Error != nil {
		return nil, result.Error
	}

	entities := make([]entities.Template, 0, len(templates))
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
	err := m.validate.Struct(template)
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

func (m mySqlTemplateRepository) Update(template entities.Template) (entities.Template, error) {
	err := m.validate.Struct(template)
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
