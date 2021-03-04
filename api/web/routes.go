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
	r.Mount("/auth", authRouter(server))
	return r
}

func unAuthRouter(server *apiServer) http.Handler {
	r := chi.NewRouter()
	r.Post("/login", server.loginHandler)
	r.Post("/register", server.registerHandler)
	return r
}

func authRouter(server *apiServer) http.Handler {
	r := chi.NewRouter()
	r.Use(ValidateJwtMiddleware)
	r.Get("/test", server.test)
	return r
}