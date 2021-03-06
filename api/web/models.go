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

type loginResponse struct {
	Jwt string `json:"jwt"`
	Username string `json:"username"`
	Role string `json:"role"`
}

type errorResponse struct {
	Message string	`json:"message"`
	StatusCode int `json:"statusCode"`
}