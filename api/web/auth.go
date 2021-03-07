package web

import (
	"github.com/go-chi/render"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	"net/http"
)

func (s *apiServer) registerHandler (w http.ResponseWriter, r *http.Request) {
	var model registrationModel
	err := render.DecodeJSON(r.Body, &model)
	if err != nil {
		createErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if model.Password != model.ConfirmPassword {
		createErrorResponse(w, "Password does not match.", http.StatusBadRequest)
		return
	}

	user := entities.User{
		Username: model.Username,
		Password: model.Password,
		RoleID:   entities.NormalRoleId,
	}

	u, err := s.userService.Create(r.Context(), user)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	token, err := tokenHandler.CreateToken(int(u.ID), u.Role.Name)
	if err != nil {
		createErrorResponse(w, "couldn't create jwt", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, loginResponse{
		Jwt: token,
		Username: u.Username,
		Role: u.Role.Name,
	})
}

func (s *apiServer) loginHandler (w http.ResponseWriter, r *http.Request) {
	var model loginCredentials
	err := render.DecodeJSON(r.Body, &model)
	if err != nil {
		createErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := s.userService.ValidateCredentials(r.Context(), model.Username, model.Password)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	token, err := tokenHandler.CreateToken(int(u.ID), u.Role.Name)
	if err != nil {
		createErrorResponse(w, "couldn't create jwt", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, loginResponse{
		Jwt: token,
		Username: u.Username,
		Role: u.Role.Name,
	})
}
