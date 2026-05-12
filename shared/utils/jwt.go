package utils

import (
	"errors"
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

// InitJWT, main.go tarafında env yüklendikten sonra çağrılmalıdır.
func InitJWT() error {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		return errors.New("JWT_SECRET environment variable is not set")
	}
	secretKey = []byte(key)
	return nil
}

//var secretKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func GenerateShortID(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	rand.Seed(time.Now().UnixNano())
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
