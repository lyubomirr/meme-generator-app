package entities

type RoleName string

const (
	AdminRoleId = 1
	NormalRoleId = 2
)

const (
	AdminRoleName RoleName = "Administrator"
	NormalRoleName RoleName = "Normal"
)

type Role struct {
	ID uint
	Name RoleName
}

type User struct {
	ID       uint
	Username string
	Password string
	Role     Role
	PictureURL string
	Memes []Meme
}

type Meme struct {
	ID uint
	Author User
	Title string
	FilePath string
	Comments []Comment
}

type Comment struct {
	ID uint
	Author User
	MemeID uint
	Content string
}

type Template struct {
	ID uint
	Name string
	FilePath string
}