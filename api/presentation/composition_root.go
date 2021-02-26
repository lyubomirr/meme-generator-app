package presentation

import (
	"github.com/lyubomirr/meme-generator/core/services"
	"github.com/lyubomirr/meme-generator/persistence"
)

func createServer() *apiServer {
	userRepo := persistence.NewUserRepository()
	authService := services.NewAuthService(userRepo)

	return &apiServer{authService: &authService}
}