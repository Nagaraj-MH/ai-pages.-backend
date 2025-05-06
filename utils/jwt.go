package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

var secretKey = []byte(os.Getenv("JWT_KEY"))

func GenerateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if email, ok := token.Claims.(jwt.MapClaims)["email"].(string); ok && token.Valid {
		return email, nil
	}

	return "", errors.New("invalid token")
}
