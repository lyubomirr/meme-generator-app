package services

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	customErr "github.com/lyubomirr/meme-generator-app/core/errors"
	"net/http"
	"strconv"
)

const (
	UserClaimsKey = "UserClaims"
	RoleClaim     = "Role"
	SubClaim      = "sub"
)

func IsAdministrator(ctx context.Context) bool {
	claims, ok := ctx.Value(UserClaimsKey).(jwt.MapClaims)
	if !ok {
		return false
	}

	r, ok := claims[RoleClaim]
	if !ok {
		return false
	}
	return r == entities.AdminRoleName
}

func getUserId(ctx context.Context) (uint, error) {
	claims, ok := ctx.Value(UserClaimsKey).(jwt.MapClaims)
	if !ok {
		return 0, customErr.NewAuthError(errors.New("couldn't get user"))
	}

	sub, ok := claims[SubClaim].(string)
	if !ok {
		return 0, customErr.NewAuthError(errors.New("couldn't get user id"))
	}

	s, err := strconv.Atoi(sub)
	if err != nil {
		return 0, customErr.NewAuthError(errors.New("couldn't get user id"))
	}

	return uint(s), nil
}

func getMimeType(file []byte) (string, error) {
	mimeType := http.DetectContentType(file)
	if !containsKey(allowedMimeTypes, mimeType) {
		return "", customErr.NewValidationError(errors.New("file content type not supported"))
	}
	return mimeType, nil
}

func containsKey(s map[string]string, e string) bool {
	for k, _ := range s {
		if k == e {
			return true
		}
	}
	return false
}

func getFileExtension(mimeType string) string {
	ext, ok := allowedMimeTypes[mimeType]
	if !ok {
		return ""
	}
	return ext
}