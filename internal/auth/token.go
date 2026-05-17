package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearerPrefix := "Bearer "
	auth := headers.Get("Authorization")
	if !strings.HasPrefix(auth, bearerPrefix) {
		return "", errors.New("Invalid/malformed authorization header")
	}

	token := auth[len(bearerPrefix):]
	return token, nil
}

func CreateJWT(
	id uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "tldr",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   id.String(),
	})

	signed, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ValidateJWT(tokenString, tokenSecret string) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{Issuer: "tldr"},
		func(t *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		},
		jwt.WithIssuer("tldr"))
	if err != nil {
		return nil, err
	}

	return token.Claims, nil
}

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(bytes)
	return token, nil
}
