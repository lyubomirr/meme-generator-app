package persistence

import "github.com/lyubomirr/meme-generator-app/core/entities"

type dbComment struct {
	ID       uint
	AuthorID uint
	Author   dbUser
	MemeID   uint
	Content  string `gorm:"type:varchar(300)"`
}

func (dbComment) TableName() string {
	return "comments"
}

func (c dbComment) toEntity() entities.Comment {
	return entities.Comment{
		ID:      c.ID,
		Author:  c.Author.toEntity(),
		Content: c.Content,
		MemeID:  c.MemeID,
	}
}

func commentsToDbModels(comments []entities.Comment) []dbComment {
	dbModels := make([]dbComment, 0, len(comments))
	for _,c := range comments {
		dbModels = append(dbModels, dbComment{
			ID:       c.ID,
			AuthorID: c.Author.ID,
			MemeID:   c.MemeID,
			Content:  c.Content,
		})
	}
	return dbModels
}