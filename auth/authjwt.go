package auth

import (
	"chat-service/config"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(roomID uint64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"room_id": roomID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"salt":    config.GlobalConfig.JWTSalt,
	})
	tokenString, err := token.SignedString(config.GlobalConfig.JWTSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (uint64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sign method: %v", token.Header["alg"])
		}
		return config.GlobalConfig.JWTSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if salt, ok := claims["salt"].(string); !ok || salt != config.GlobalConfig.JWTSalt {
			return 0, errors.New("invalid salt")
		}

		if roomIDFloat, ok := claims["room_id"].(float64); ok {
			return uint64(roomIDFloat), nil
		}

		return 0, errors.New("invalid room_id")
	}

	return 0, errors.New("invalid token")
}
