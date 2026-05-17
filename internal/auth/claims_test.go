package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGetUserID_ValidClaims(t *testing.T) {
	claims := &jwt.RegisteredClaims{
		Subject: uuid.Must(uuid.NewRandom()).String(),
	}
	expected := claims.Subject
	ctx := context.WithValue(context.Background(), ClaimsKey, claims)

	userId, err := GetUserID(ctx)
	if err != nil {
		t.Fatalf("got %q", err.Error())
	}

	if userId.String() != expected {
		t.Fatalf("want %q, got %q", expected, userId.String())
	}
}

func TestGetUserID_InvalidClaims(t *testing.T) {
	claims := &jwt.MapClaims{
		"subject": uuid.Must(uuid.NewRandom()).String(),
	}
	expected := errors.New("invalid claims")
	ctx := context.WithValue(context.Background(), ClaimsKey, claims)

	_, err := GetUserID(ctx)
	if err.Error() != expected.Error() {
		t.Fatalf("want %q, got %q", expected.Error(), err.Error())
	}
}

func TestGetUserID_InvalidSubject(t *testing.T) {
	claims := &jwt.RegisteredClaims{
		Subject: "123456",
	}
	ctx := context.WithValue(context.Background(), ClaimsKey, claims)

	_, err := GetUserID(ctx)
	if err == nil {
		t.Fatalf("expected err to be non-nil")
	}
}
