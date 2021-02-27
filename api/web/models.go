package web

import "github.com/lyubomirr/meme-generator-app/core/entities"

type loginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registrationModel struct {
	entities.User
	ConfirmPassword string `json: "confirmPassword"`
}

type jwtResponse struct {
	Jwt string `json:"jwt"`
}
