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

	r.Route("/api", func(r chi.Router) {
		r.Mount("/", unAuthRouter(server))
		r.Mount("/auth", authRouter(server))
	})
	return r
}

func unAuthRouter(server *apiServer) http.Handler {
	r := chi.NewRouter()
	r.Post("/login", server.loginHandler)
	r.Post("/register", server.registerHandler)
	r.Get("/template/file/{id}", server.getTemplateFileHandler)
	r.Get("/meme/file/{id}", server.getMemeFileHandler)

	return r
}

func authRouter(server *apiServer) http.Handler {
	r := chi.NewRouter()
	r.Use(validateJwtMiddleware)
	r.Mount("/admin", adminRouter(server))

	r.Route("/template", func(r chi.Router) {
		r.Get("/", server.getTemplatesHandler)
		r.Get("/{id}", server.getTemplateHandler)
	})

	r.Route("/meme", func(r chi.Router) {
		r.Get("/", server.getMemesHandler)
		r.Post("/", server.addMemeHandler)
		r.Get("/{id}", server.getMemeHandler)
		r.Delete("/{id}", server.deleteMemeHandler)
		r.Post("/{id}/comment", server.addCommentHandler)
		r.Delete("/{memeId}/comment/{commentId}", server.deleteCommentHandler)
	})

	return r
}

func adminRouter(server *apiServer) http.Handler {
	r := chi.NewRouter()
	r.Use(adminOnlyMiddleware)

	r.Route("/template", func(r chi.Router) {
		r.Post("/", server.addTemplateHandler)
		r.Delete("/{id}", server.deleteTemplateHandler)
	})
	return r
}