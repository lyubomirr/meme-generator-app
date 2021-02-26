package web

import (
	"github.com/lyubomirr/meme-generator-app/core/services"
	"log"
	"net/http"
)

type apiServer struct {
	authService *services.Authentication
}

func Serve(address string) {
	server := createServer()
	router := constructRouter(server)
	log.Fatal(http.ListenAndServe(address, router))
}