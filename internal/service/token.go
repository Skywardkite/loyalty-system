package service

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("my_super_secret_goose")

func GenerateAuthToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // токен действителен 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен супер секретным ключом
	return token.SignedString(jwtSecret)
}

func ParseAuthToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenInvalidId
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, jwt.ErrTokenInvalidClaims
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, jwt.ErrTokenInvalidClaims
	}

	return int64(userIDFloat), nil
}
