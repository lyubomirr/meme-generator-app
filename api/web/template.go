package web

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strconv"
)

func (s *apiServer) addTemplateHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := parseMultipartFile(r, "file")
	if err != nil {
		createErrorResponse(w, fmt.Sprintf("cannot read template file: %v", err.Error()), http.StatusBadRequest)
		return
	}

	templateJson := r.FormValue("template")
	if templateJson == "" {
		createErrorResponse(w, "no template object sent", http.StatusBadRequest)
		return
	}

	var template entities.Template
	err = json.Unmarshal([]byte(templateJson), &template)
	if err != nil {
		createErrorResponse(w, fmt.Sprintf("cannot unmarshal template: %v", err.Error()), http.StatusBadRequest)
		return
	}

	template, err = s.templateService.Create(r.Context(), bytes, template)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	render.JSON(w, r, template)
}

func (s *apiServer) deleteTemplateHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		createErrorResponse(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = s.templateService.Delete(r.Context(), uint(id))
	if err != nil {
		handleServiceError(w, err)
	}
}

func (s *apiServer) getTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	templates, err := s.templateService.GetAll(r.Context())
	if err != nil {
		handleServiceError(w, err)
		return
	}
	render.JSON(w, r, templates)
}

func (s *apiServer) getTemplateHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		createErrorResponse(w, "invalid id", http.StatusBadRequest)
		return
	}

	t, err := s.templateService.Get(r.Context(), uint(id))
	if err != nil {
		handleServiceError(w, err)
		return
	}
	render.JSON(w, r, t)
}

func (s *apiServer) getTemplateFileHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		createErrorResponse(w, "invalid id", http.StatusBadRequest)
		return
	}

	t, err := s.templateService.Get(r.Context(), uint(id))
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Disposition", "inline")
	http.ServeFile(w, r, t.FilePath)
}