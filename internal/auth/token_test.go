package auth

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGetBearerToken_ValidHeader(t *testing.T) {
	headers := http.Header{"Authorization": []string{"Bearer 123456"}}
	expected := "123456"
	actual, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("expected err to be nil, got %q", err.Error())
	}

	if expected != actual {
		t.Fatalf("want %q, got %q", expected, actual)
	}
}

func TestBearerToken_InvalidHeader(t *testing.T) {
	var tests = []struct {
		name   string
		header http.Header
	}{
		{
			name:   "Missing authorization header",
			header: http.Header{},
		},
		{
			name:   "Missing authorization value",
			header: http.Header{"Authorization": []string{}},
		},
		{
			name:   "Missing bearer token",
			header: http.Header{"Authorization": []string{"Bearer"}},
		},
		{
			name:   "Invalid authorization format",
			header: http.Header{"Authorization": []string{"Credentials 123"}},
		},
	}

	expected := errors.New("Invalid/malformed authorization header")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBearerToken(tt.header)
			if err == nil {
				t.Fatal("expected error to be non-nil")
			}

			if err.Error() != expected.Error() {
				t.Fatalf("want %q, got %q", expected.Error(), err.Error())
			}
		})
	}
}

func TestCreateJWT(t *testing.T) {
	id := uuid.Must(uuid.NewRandom())
	secret := "SeCrEtKeY1!2@3#"
	expiresIn := time.Minute * 10
	iss := "tldr"

	tokenString, err := CreateJWT(id, secret, expiresIn)
	if err != nil {
		t.Fatalf("expected error to be nil, got %q", err.Error())
	}

	if strings.TrimSpace(tokenString) == "" {
		t.Fatal("expected token string to be non-empty")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{Issuer: iss},
		func(t *jwt.Token) (any, error) { return []byte(secret), nil },
		jwt.WithIssuer(iss),
	)
	if err != nil {
		t.Fatalf("expected error to be nil, got %q", err.Error())
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		t.Fatalf("expected claims to be of type *jwt.RegisteredClaims")
	}

	if claims.Issuer != iss {
		t.Fatalf("want %q, got %q", iss, claims.Issuer)
	}

	if claims.Subject != id.String() {
		t.Fatalf("want %q, got %q", id.String(), claims.Subject)
	}

	until := time.Until(claims.ExpiresAt.Time)

	diff := until - expiresIn
	if diff < 0 {
		diff = -diff
	}

	if diff > time.Second {
		t.Fatalf(
			"expiry diff too large: got %v, want within %v of %v",
			until,
			time.Second,
			expiresIn,
		)
	}
}

func TestValidateJWT_ValidToken(t *testing.T) {
	id := uuid.Must(uuid.NewRandom())
	secret := "SeCrEtKeY1!2@3#"
	expiresIn := time.Minute * 10
	iss := "tldr"

	tokenString, err := CreateJWT(id, secret, expiresIn)
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err.Error())
	}

	claims, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err.Error())
	}

	registeredClaims, ok := claims.(*jwt.RegisteredClaims)
	if !ok {
		t.Fatalf("expected claims to be of type *jwt.RegisteredClaims")
	}

	if registeredClaims.Issuer != iss {
		t.Fatalf("want %q, got %q", iss, registeredClaims.Issuer)
	}

	if registeredClaims.Subject != id.String() {
		t.Fatalf("want %q, got %q", id.String(), registeredClaims.Subject)
	}

	until := time.Until(registeredClaims.ExpiresAt.Time)

	diff := until - expiresIn
	if diff < 0 {
		diff = -diff
	}

	if diff > time.Second {
		t.Fatalf(
			"expiry diff too large: got %v, want within %v of %v",
			until,
			time.Second,
			expiresIn,
		)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	secret := "SeCrEtKeY1!2@3#"
	_, err := ValidateJWT("invalidtoken", secret)
	if err == nil {
		t.Fatal("expected error to be non-nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	id := uuid.Must(uuid.NewRandom())
	secret := "SeCrEtKeY1!2@3#"
	expiresIn := time.Minute * 10

	tokenString, err := CreateJWT(id, secret, expiresIn)
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err.Error())
	}
	_, err = ValidateJWT(tokenString, "wrongsecret")
	if err == nil {
		t.Fatal("expected error to be non-nil")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	id := uuid.Must(uuid.NewRandom())
	secret := "SeCrEtKeY1!2@3#"
	expiresIn := -time.Second

	tokenString, err := CreateJWT(id, secret, expiresIn)
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err.Error())
	}
	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Fatal("expected error to be non-nil")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	const expectedLength = 64
	token, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err.Error())
	}

	if len(token) != expectedLength {
		t.Fatalf(
			"unexpected token length: want %v, got %v",
			expectedLength,
			len(token),
		)
	}
}
