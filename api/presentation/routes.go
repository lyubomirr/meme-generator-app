package presentation

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func GetRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/", unAuthRouter())
	return r
}

func unAuthRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/login", server.loginHandler)
	r.Get("/register", server.registerHandler)
	return r
}