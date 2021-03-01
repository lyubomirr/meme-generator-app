package services

import (
	"context"
	"github.com/lyubomirr/meme-generator-app/core/entities"
)

type Meme interface {
	GetAll(ctx context.Context) ([]entities.Meme, error)
	Get(ctx context.Context, id uint) (entities.Meme, error)
	GetByAuthor(ctx context.Context, userId uint) ([]entities.Meme, error)
}