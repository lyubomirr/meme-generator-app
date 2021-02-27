package web

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/lyubomirr/meme-generator-app/core/entities"
	"strconv"
	"time"
)

type jwtHandler struct {
	Secret string
	Lifetime time.Duration
	Issuer string
	Audience string
}

type claims struct {
	Role string
	jwt.StandardClaims
}

func (j *jwtHandler) CreateToken(userId int, role entities.RoleName) (string , error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		StandardClaims: jwt.StandardClaims{
			Audience:  j.Audience,
			ExpiresAt: time.Now().Local().Add(j.Lifetime).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    j.Issuer,
			NotBefore: time.Now().Unix(),
			Subject:   strconv.Itoa(userId),
		},
		Role: string(role),
	})

	signed, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

