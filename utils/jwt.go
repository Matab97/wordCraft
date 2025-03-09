package utils

import (
	"errors"
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "secret"

func GenerateJWT(userId string, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (string, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	tokenIsInvalid := parsedToken.Valid
	if !tokenIsInvalid {
		return "", errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token")
	}
	if claims["exp"] == nil || claims["exp"].(float64) < float64(time.Now().Unix()) {
		return "", errors.New("invalid token")
	}
	// email := claims["email"].(string)
	userId := claims["userId"].(string)
	return userId, nil
}
