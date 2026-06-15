package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

func TestTokenWithIncorrectSignString(t *testing.T) {
	t.Setenv("JWT_SECRET", "incorrent-secret")

	incorrectToken, err := GenerateJwtToken("test@mail.com", 1)

	if incorrectToken == "" || err != nil {
		t.Fatalf("Token generation failed. token: %v, error: %v", incorrectToken, err)
	}

	t.Setenv("JWT_SECRET", "corrent-secret")

	claim, err := VerifyToken(incorrectToken)

	if !errors.Is(err, jwt.ErrTokenSignatureInvalid) || len(claim) != 0 {
		t.Errorf("Token had been succefulery verified. claim: %v, error: %v", claim, err)
	}
}

func TestTokenGeneration(t *testing.T) {
	t.Parallel()

	hashedPassword, err := GenerateJwtToken("test@email.com", 1)

	if err != nil {
		t.Fatalf("Token generation failed. error: %v", err)
	}

	claim, err := VerifyToken(hashedPassword)

	if err != nil {
		t.Fatalf("Vefication token is failed. error: %v", err)
	}

	if claim["userEmail"] != "test@email.com" || int64(claim["userId"].(float64)) != 1 {
		t.Errorf("Claim does not contain correct data. claim: %v", claim)
	}
}

func TestTokenWithDiffEncriptMethod(t *testing.T) {
	t.Parallel()

	token := jwt.New(jwt.SigningMethodEdDSA)

	_, privKey, err := ed25519.GenerateKey(rand.Reader)
	ed25519.Sign(privKey, []byte("test-secrete"))

	if err != nil {
		t.Fatalf("Failed to generate sing key. error: %v", err)
	}

	strToken, err := token.SignedString(privKey)

	if err != nil {
		t.Fatalf("Failed to generate token. error: %v", err)
	}

	claim, err := VerifyToken(strToken)

	if !errors.Is(err, jwt.ErrTokenSignatureInvalid) || len(claim) > 0 {
		t.Errorf("Token was succesfuly verifyied or other error happend. error: %v, claim: %v", err, claim)
	}
}

func TestTokenWithoutClaim(t *testing.T) {
	t.Parallel()

	token := jwt.New(jwt.SigningMethodHS256)
	strToken, err := token.SignedString(getJwtSecret())

	if err != nil {
		t.Fatalf("Token generation failed. error: %v", err)
	}

	claim, err := VerifyToken(strToken)

	if len(claim) > 0 || err.Error() != "Invalid Token" {
		t.Errorf("Token was succesfuly verifyied or other error happend. error: %v, claim: %v", err, claim)
	}
}

func TestTokenWithIncorrectFormatForExp(t *testing.T) {
	t.Parallel()

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{"exp": "Incorrect exp"},
	)

	strToken, err := token.SignedString(getJwtSecret())

	if err != nil {
		t.Fatalf("Token generation failed. error: %v", err)
	}

	claim, err := VerifyToken(strToken)

	if len(claim) > 0 || !errors.Is(err, jwt.ErrInvalidType) {
		t.Errorf("Token was succesfuly verifyied or other error happend. error: %v, claim: %v", err, claim)
	}
}

func TestExpieredToken(t *testing.T) {
	t.Parallel()

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{"exp": time.Now().Unix()},
	)

	strToken, err := token.SignedString(getJwtSecret())

	if err != nil {
		t.Fatalf("Token generation failed. error: %v", err)
	}

	claim, err := VerifyToken(strToken)

	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Errorf("Test. claim: %v, error: %v", claim, err)
	}
}

func TestVerifyTokenWithRandomString(t *testing.T) {
	t.Parallel()

	claim, err := VerifyToken("some-random-string")

	if !errors.Is(err, jwt.ErrTokenMalformed) {
		t.Errorf("Test. claim: %v, error: %v", claim, err)
	}
}
