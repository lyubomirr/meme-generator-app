package web

import (
	"github.com/lyubomirr/meme-generator-app/core/services"
	"github.com/lyubomirr/meme-generator-app/persistence"
	"os"
	"time"
)

func initJwtHandler() *jwtHandler {
	return &jwtHandler{
		Secret:   os.Getenv("jwt_secret"),
		Lifetime: time.Hour * 24,
		Issuer:   "meme-generator-auth",
		Audience: "meme-generator-app",
	}
}

func createServer() *apiServer {
	uowFactory := persistence.NewUnitOfWorkFactory()
	userService := services.NewUserService(uowFactory)
	memeService := services.NewMemeService(uowFactory)
	templateService := services.NewTemplateService(uowFactory)

	return &apiServer{
		userService:     userService,
		memeService:     memeService,
		templateService: templateService,
	}
}
