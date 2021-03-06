package web

import (
	"github.com/lyubomirr/meme-generator-app/core/services"
	"log"
	"net/http"
)

type apiServer struct {
	userService     services.User
	memeService     services.Meme
	templateService services.Template
}

func Serve(address string) {
	tokenHandler = initJwtHandler()
	server := createServer()
	router := constructRouter(server)
	log.Fatal(http.ListenAndServe(address, router))
}