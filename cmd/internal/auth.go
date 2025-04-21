package internal

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var JwtSecret = []byte(os.Getenv("SECRET_KEY"))

func GenerateToken(username string, email string) (string, error) {
	now := time.Now().UTC()

	claims := jwt.MapClaims{
		"username": username,
		"email":    email,
		"role":     "",
		"iat":      now.Unix(),
		"nbf":      now.Unix(),
		"exp":      now.Add(time.Hour * 24).Unix(), // Token validity is 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", errors.New("unable to signed the token")
	}

	return tokenStr, nil
}

func ValidateToken(tokenString string) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}
	return claims["email"], nil // TODO return the correct information
}
