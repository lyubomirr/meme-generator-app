package web

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	"net/http"
	"strconv"
)

func (s *apiServer) addMemeHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := parseMultipartFile(r, "file")
	if err != nil {
		createErrorResponse(w, fmt.Sprintf("cannot read meme file: %v", err.Error()), http.StatusBadRequest)
		return
	}

	memeJson := r.FormValue("meme")
	if memeJson == "" {
		createErrorResponse(w, "no meme object sent", http.StatusBadRequest)
		return
	}

	var meme entities.Meme
	err = json.Unmarshal([]byte(memeJson), &meme)
	if err != nil {
		createErrorResponse(w, fmt.Sprintf("cannot unmarshal meme: %v", err.Error()), http.StatusBadRequest)
		return
	}

	meme, err = s.memeService.Create(r.Context(), bytes, meme)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	render.JSON(w, r, meme)
}

func (s *apiServer) deleteMemeHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		createErrorResponse(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = s.memeService.Delete(r.Context(), uint(id))
	if err != nil {
		handleServiceError(w, err)
	}
}

func (s *apiServer) getMemesHandler(w http.ResponseWriter, r *http.Request) {
	memes, err := s.memeService.GetAll(r.Context())
	if err != nil {
		handleServiceError(w, err)
		return
	}
	render.JSON(w, r, memes)
}

func (s *apiServer) getMemeHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		createErrorResponse(w, "invalid id", http.StatusBadRequest)
		return
	}

	meme, err := s.memeService.Get(r.Context(), uint(id))
	if err != nil {
		handleServiceError(w, err)
		return
	}
	render.JSON(w, r, meme)
}

func (s *apiServer) getMemeFileHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		createErrorResponse(w, "invalid id", http.StatusBadRequest)
		return
	}

	m, err := s.memeService.Get(r.Context(), uint(id))
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Disposition", "inline")
	http.ServeFile(w, r, m.FilePath)
}

func (s *apiServer) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		createErrorResponse(w, "invalid id", http.StatusBadRequest)
		return
	}

	var comment entities.Comment
	err = render.DecodeJSON(r.Body, &comment)
	if err != nil {
		createErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := s.memeService.AddComment(r.Context(), uint(id), comment)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	render.JSON(w, r, m)
}

func (s *apiServer) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	memeIdParam := chi.URLParam(r, "memeId")
	memeId, err := strconv.Atoi(memeIdParam)
	if err != nil {
		createErrorResponse(w, "invalid meme id", http.StatusBadRequest)
		return
	}

	commentIdParam := chi.URLParam(r, "commentId")
	commentId, err := strconv.Atoi(commentIdParam)
	if err != nil {
		createErrorResponse(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	meme, err := s.memeService.DeleteComment(r.Context(), uint(memeId), uint(commentId))
	if err != nil {
		handleServiceError(w, err)
		return
	}
	render.JSON(w, r, meme)
}