package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken creates a JWT token for a user.
func GenerateToken(userID, key string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}