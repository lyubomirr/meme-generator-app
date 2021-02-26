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
	r.Get("/login", server.loginHandler)
	r.Get("/register", server.registerHandler)
	return r
}