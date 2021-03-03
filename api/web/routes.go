package web

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func constructRouter(server *apiServer) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/", unAuthRouter(server))
	return r
}

func unAuthRouter(server *apiServer) http.Handler {
	r := chi.NewRouter()
	r.Post("/login", server.loginHandler)
	r.Post("/register", server.registerHandler)
	r.Post("/test", server.test)
	return r
}