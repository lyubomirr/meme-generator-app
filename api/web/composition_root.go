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
	userService := services.NewUserService(uowFactory)
	memeService := services.NewMemeService(uowFactory)
	templateService := services.NewTemplateService(uowFactory)

	return &apiServer{
		userService:     userService,
		memeService:     memeService,
		templateService: templateService,
	}
}
