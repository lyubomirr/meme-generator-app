package web

import (
	"encoding/json"
	customErr "github.com/lyubomirr/meme-generator-app/core/errors"
	"io/ioutil"
	"net/http"
)

func handleServiceError(w http.ResponseWriter, err error) {
	switch v := err.(type) {
	case customErr.AuthError:
		createErrorResponse(w, v.Error(), http.StatusUnauthorized)
	case customErr.ValidationError:
		createErrorResponse(w, v.Error(), http.StatusBadRequest)
	case customErr.RightsError:
		createErrorResponse(w, v.Error(), http.StatusForbidden)
	default:
		createErrorResponse(w, v.Error(), http.StatusInternalServerError)
	}
}

func createErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	resp := errorResponse{
		Message:    message,
		StatusCode: statusCode,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "", statusCode)
		return
	}

	http.Error(w, string(jsonResp), statusCode)
	return
}

func parseMultipartFile(r *http.Request, key string) ([]byte, error) {
	r.ParseMultipartForm((1 << 20) * 10)
	file, _, err := r.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

