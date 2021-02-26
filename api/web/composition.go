package web

import (
	"github.com/lyubomirr/meme-generator-app/core/services"
	"github.com/lyubomirr/meme-generator-app/persistence"
)

func createServer() *apiServer {
	userRepo := persistence.NewUserRepository()
	authService := services.NewAuthService(userRepo)

	return &apiServer{authService: &authService}
}