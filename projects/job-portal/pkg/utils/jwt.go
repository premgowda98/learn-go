package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(username string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	token, err := claims.SignedString([]byte("secret"))

	if err != nil {
		return "", err
	}

	return token, nil

}
