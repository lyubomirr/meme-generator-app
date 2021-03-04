package web

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type jwtHandler struct {
	Secret string
	Lifetime time.Duration
	Issuer string
	Audience string
}

var tokenHandler *jwtHandler

type claims struct {
	Role string
	jwt.StandardClaims
}

func (j *jwtHandler) CreateToken(userId int, role string) (string , error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		StandardClaims: jwt.StandardClaims{
			Audience:  j.Audience,
			ExpiresAt: time.Now().Local().Add(j.Lifetime).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    j.Issuer,
			NotBefore: time.Now().Unix(),
			Subject:   strconv.Itoa(userId),
		},
		Role: role,
	})

	signed, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (j *jwtHandler) ValidateToken(token string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(jwt *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if err != nil {
		return jwt.MapClaims{}, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid jwt")
	}
}


