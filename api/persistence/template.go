package persistence

import (
	"github.com/go-playground/validator/v10"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type dbTemplateTextbox struct {
	TemplateID uint    `gorm:"primaryKey"`
	TopOffset  float64 `gorm:"primaryKey"`
	LeftOffset float64 `gorm:"primaryKey"`
	Width      float64
	Height     float64
}

func (dbTemplateTextbox) TableName() string {
	return "template_textboxes"
}

type dbTemplate struct {
	ID        uint
	Name      string `gorm:"type:varchar(50)"`
	FilePath  string
	MimeType  string              `gorm:"type:varchar(50)"`
	Textboxes []dbTemplateTextbox `gorm:"foreignKey:TemplateID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time
}

func (dbTemplate) TableName() string {
	return "templates"
}

func textboxesToEntities(textboxes []dbTemplateTextbox) []entities.TemplateTextbox {
	e := make([]entities.TemplateTextbox, 0, len(textboxes))
	for _, t := range textboxes {
		e = append(e, entities.TemplateTextbox{
			TopOffset:  t.TopOffset,
			LeftOffset: t.LeftOffset,
			Width:      t.Width,
			Height:     t.Height,
		})
	}
	return e
}

func newTextboxes(textboxes []entities.TemplateTextbox) []dbTemplateTextbox {
	dbModels := make([]dbTemplateTextbox, 0, len(textboxes))
	for _, t := range textboxes {
		dbModels = append(dbModels, dbTemplateTextbox{
			LeftOffset: t.LeftOffset,
			TopOffset:  t.TopOffset,
			Width:      t.Width,
			Height:     t.Height,
		})
	}
	return dbModels
}

func (m dbTemplate) toEntity() entities.Template {
	return entities.Template{
		ID:        m.ID,
		Name:      m.Name,
		FilePath:  m.FilePath,
		MimeType:  m.MimeType,
		Textboxes: textboxesToEntities(m.Textboxes),
		CreatedAt: m.CreatedAt,
	}
}

func newTemplate(template entities.Template) dbTemplate {
	return dbTemplate{
		ID:        template.ID,
		Name:      template.Name,
		FilePath:  template.FilePath,
		MimeType:  template.MimeType,
		Textboxes: newTextboxes(template.Textboxes),
		CreatedAt: template.CreatedAt,
	}
}

type mySqlTemplateRepository struct {
	db       *gorm.DB
	validate *validator.Validate
}

func (m mySqlTemplateRepository) GetAll() ([]entities.Template, error) {
	var templates []dbTemplate
	result := withReadTimeout(m.db).
		Preload(clause.Associations).
		Order("created_at desc").
		Find(&templates)

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
	result = m.db.Save(&dbTemplate)
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
