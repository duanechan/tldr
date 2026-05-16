package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const ClaimsKey = "claims"

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	claims, ok := ctx.Value(ClaimsKey).(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid claims")
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}
