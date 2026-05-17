package utils

import (
	//"fmt"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateJwtToken(email string, id int64) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userEmail": email,
			"userId":    id,
			"exp":       time.Now().Add(time.Hour * 2).Unix(),
		},
	)

	return token.SignedString(getJwtSecret())
}

func VerifyToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(
		token,
		func(token *jwt.Token) (any, error) { return getJwtSecret(), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return nil, err
	}

	exp, ok := claims["exp"].(float64)

	if !ok {
		return nil, errors.New("Invalid Token")
	}

	expDate := time.Unix(int64(exp), 0)

	if !time.Now().Before(expDate) {
		return nil, errors.New("Invalid Token")
	}

	return claims, err
}

func getJwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}
