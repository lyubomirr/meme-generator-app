package entities

import "time"

type RoleName string

const (
	AdminRoleId  = 1
	NormalRoleId = 2
)

const (
	AdminRoleName  = "Administrator"
	NormalRoleName = "Normal"
)

type Role struct {
	ID   uint `validate:"required"`
	Name string
}

type User struct {
	ID         uint   `json:"id"`
	Username   string `json:"username" validate:"required,max=25"`
	Password   string `json:"-" validate:"required,min=8"`
	RoleID     uint   `json:"roleId" validate:"required"`
	Role       Role   `json:"role" validate:"-"`
	Memes      []Meme `json:"memes,omitempty"`
}

type Meme struct {
	ID         uint      `json:"id"`
	AuthorID   uint      `json:"authorId" validate:"required"`
	Author     User      `json:"author" validate:"-"`
	Title      string    `json:"title" validate:"required,max=50"`
	FilePath   string    `json:"-" validate:"required"`
	MimeType   string    `json:"mimeType" validate:"required,max=50`
	CreatedAt  time.Time `json:"createdAt"`
	Comments   []Comment `json:"comments"`
	TemplateID uint      `json:"templateId"`
	Template   Template  `json:"template" validate:"-"`
}

type Comment struct {
	ID        uint      `json:"id"`
	AuthorID  uint      `json:"authorId" validate:"required"`
	Author    User      `json:"author" validate:"-"`
	MemeID    uint      `json:"memeId" validate:"required"`
	Content   string    `json:"content" validate:"required,max=50"`
	CreatedAt time.Time `json:"createdAt"`
}

type Template struct {
	ID            uint                   `json:"id"`
	Name          string                 `json:"name" validate:"required,max=50"`
	FilePath      string                 `json:"-" validate:"required"`
	MimeType      string                 `json:"mimeType" validate:"required,max=50`
	TextPositions []TemplateTextPosition `json:"textPositions" validate:"required"`
	CreatedAt     time.Time
}

type TemplateTextPosition struct {
	TopOffset  uint `json:"topOffset"`
	LeftOffset uint `json:"leftOffset"`
}
