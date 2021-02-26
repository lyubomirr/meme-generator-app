package presentation

import (
	"github.com/lyubomirr/meme-generator/core/services"
)

type apiServer struct {
	authService *services.Authentication
}

var server *apiServer

func init() {
	server = createServer()
}