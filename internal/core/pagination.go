package core

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func extractQueryParams(query url.Values) (
	*time.Time,
	uuid.UUID,
	int64,
	[]FieldError,
) {
	var fieldErrors []FieldError

	cursor := strings.TrimSpace(query.Get("cursor"))
	if cursor == "" {
		cursor = DefaultPageCursor
	}

	createdAt, id, err := decodeCursor(cursor)
	if err != nil {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   "cursor",
			Message: "Invalid 'cursor' format",
		})
	}

	limitStr := strings.TrimSpace(query.Get("limit"))
	if limitStr == "" {
		limitStr = DefaultPageLimit
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   "limit",
			Message: "Invalid 'limit' format",
		})
	} else if limit <= 0 || limit > 100 {
		fieldErrors = append(fieldErrors, FieldError{
			Field: "limit",
			Message: fmt.Sprintf("'limit' must be between 1 and %d",
				MaxPageLimit,
			),
		})
	}

	if len(fieldErrors) > 0 {
		return nil, uuid.Nil, 0, fieldErrors
	}

	return createdAt, id, limit, nil
}

func decodeCursor(cursor string) (*time.Time, uuid.UUID, error) {
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, uuid.Nil, err
	}

	parts := strings.Split(string(raw), PageCursorSeparator)
	if len(parts) < 2 {
		return nil, uuid.Nil, errors.New("invalid cursor format")
	}

	seconds, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, uuid.Nil, err
	}

	createdAt := time.Unix(seconds, 0).UTC()

	id, err := uuid.Parse(parts[1])
	if err != nil {
		return nil, uuid.Nil, err
	}

	return &createdAt, id, nil
}

func encodeCursor(createdAt *time.Time, id uuid.UUID) PageCursor {
	unixStr := strconv.FormatInt(createdAt.Unix(), 10)
	formatted := unixStr + PageCursorSeparator + id.String()
	next := base64.RawURLEncoding.EncodeToString([]byte(formatted))
	return PageCursor(string(next))
}
