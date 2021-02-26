package persistence

import (
	"github.com/lyubomirr/meme-generator-app/core/entities"
)

type dbRole struct {
	ID uint
	Name entities.RoleName `gorm:"type:varchar(25);uniqueIndex"`
}

func (dbRole) TableName() string {
	return "roles"
}

func (r dbRole) toEntity() entities.Role {
	return entities.Role{
		ID:   r.ID,
		Name: r.Name,
	}
}