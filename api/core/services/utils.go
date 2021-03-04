package services

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/lyubomirr/meme-generator-app/core/entities"
)

const (
	UserClaimsKey = "UserClaims"
	RoleClaim     = "Role"
	SubClaim      = "sub"
)

func isAdministrator(ctx context.Context) bool {
	claims := ctx.Value(UserClaimsKey).(jwt.MapClaims)
	r, ok := claims[RoleClaim]
	if !ok {
		return false
	}
	return r == entities.AdminRoleName
}
