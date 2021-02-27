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
}

