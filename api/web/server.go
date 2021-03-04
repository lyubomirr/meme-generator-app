package web

import (
	"github.com/lyubomirr/meme-generator-app/core/services"
	"log"
	"net/http"
)

type apiServer struct {
	authService services.User
	memeService services.Meme
}

func Serve(address string) {
	tokenHandler = initJwtHandler()
	server := createServer()
	router := constructRouter(server)
	log.Fatal(http.ListenAndServe(address, router))
}