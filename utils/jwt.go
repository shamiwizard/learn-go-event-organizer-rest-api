package utils

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(email string, id int64) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userEmail": email,
			"userId": id,
			"exp": time.Now().Add(time.Hour * 2).Unix(),
		},
	)

	return token.SignedString(getJwtSecret())
}

func getJwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

