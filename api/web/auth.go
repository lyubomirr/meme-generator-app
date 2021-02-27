package web

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	customErr "github.com/lyubomirr/meme-generator-app/core/errors"
	"net/http"
)

func (s *apiServer) registerHandler (w http.ResponseWriter, r *http.Request) {
	var model registrationModel
	err := render.DecodeJSON(r.Body, &model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if model.Password != model.ConfirmPassword {
		http.Error(w, "Password does not match.", http.StatusBadRequest)
		return
	}

	model.User.Role = entities.Role{ID: entities.NormalRoleId}
	u, err := s.authService.Create(model.User)
	if err != nil {
		if errors.Is(err, customErr.ValidationError{}) || errors.Is(err, customErr.ExistingResourceError{}) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := s.jwt.CreateToken(int(u.ID), u.Role.Name)
	if err != nil {
		http.Error(w, "couldn't create jwt", http.StatusInternalServerError)
	}
	render.JSON(w, r, jwtResponse{Jwt: token})
}

func (s *apiServer) loginHandler (w http.ResponseWriter, r *http.Request) {
	var model loginCredentials
	err := render.DecodeJSON(r.Body, &model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := s.authService.ValidateCredentials(model.Username, model.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := s.jwt.CreateToken(int(u.ID), u.Role.Name)
	if err != nil {
		http.Error(w, "couldn't create jwt", http.StatusInternalServerError)
	}
	render.JSON(w, r, jwtResponse{Jwt: token})
}
