package web

import (
	"github.com/lyubomirr/meme-generator-app/core/services"
	"github.com/lyubomirr/meme-generator-app/persistence"
	"time"
)

func createServer() *apiServer {
	uowFactory := persistence.NewUnitOfWorkFactory()
	authService := services.NewUserService(uowFactory)
	jwtHandler := jwtHandler{
		Secret:   "secret12312312",
		Lifetime: time.Hour * 24,
		Issuer:   "meme-generator-auth",
		Audience: "meme-generator-app",
	}

	return &apiServer{authService: authService, jwt: &jwtHandler}
}