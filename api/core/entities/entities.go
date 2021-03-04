package entities

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
	ID         uint
	Username   string `validate:"required,max=25"`
	Password   string `validate:"required,min=8"`
	Role       Role   `validate:"required"`
	PictureURL string
	Memes      []Meme
}

type Meme struct {
	ID       uint
	Author   User   `validate:"required"`
	Title    string `validate:"required,max=50"`
	FilePath string `validate:"required"`
	MimeType string `validate:"required,max=50`
	Comments []Comment
	Template Template
}

type Comment struct {
	ID      uint
	Author  User   `validate:"required"`
	MemeID  uint   `validate:"required"`
	Content string `validate:"required,max=50"`
}

type Template struct {
	ID            uint
	Name          string                 `validate:"required,max=50"`
	FilePath      string                 `validate:"required"`
	MimeType      string                 `validate:"required,max=50`
	TextPositions []TemplateTextPosition `validate:"required"`
}

type TemplateTextPosition struct {
	TopOffset  uint
	LeftOffset uint
}
