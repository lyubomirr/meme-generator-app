package web

import (
	"github.com/lyubomirr/meme-generator-app/core/services"
	"github.com/lyubomirr/meme-generator-app/persistence"
	"time"
)

func initJwtHandler() *jwtHandler {
	return &jwtHandler{
		Secret:   "secret12312312",
		Lifetime: time.Hour * 24,
		Issuer:   "meme-generator-auth",
		Audience: "meme-generator-app",
	}
}

func createServer() *apiServer {
	uowFactory := persistence.NewUnitOfWorkFactory()
	authService := services.NewUserService(uowFactory)
	memeService := services.NewMemeService(uowFactory)

	return &apiServer{authService: authService, memeService: memeService}
}