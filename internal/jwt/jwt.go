package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// NewToken generates a new JWT token.
func NewToken(email string, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
