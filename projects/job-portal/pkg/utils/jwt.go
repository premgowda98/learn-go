package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(username string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"username": username,
		"exp": time.Now().Add(time.Hour*5).Unix(),
		"iat": time.Now().Unix(),
	})

	token, err := claims.SignedString([]byte("secret"))

	if err != nil {
		return "", err
	}

	return token, nil

}

func ValidateJWTToken(token string) (*Claims, error) {
	claims := &Claims{}

	tokenParsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil || !tokenParsed.Valid {
		return nil, err
	}

	return claims, nil
}
