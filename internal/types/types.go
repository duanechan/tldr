package types

import (
	"encoding/json"
	"time"
)

type Page[T any] struct {
	Next    *PageCursor `json:"next,omitempty"`
	Results []T         `json:"results"`
}

type PageCursor time.Time

func (p PageCursor) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(p).Format(time.RFC3339))
}

type PageLimit int64

type HealthResponse struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type SummarizeResponse struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Duration int64  `json:"duration,omitempty"`
}

type SummarizeTextRequest struct {
	Text string `json:"text"`
}

type ErrorResponse struct {
	Code      int          `json:"code"`
	RequestID string       `json:"request_id"`
	Message   string       `json:"message,omitempty"`
	Errors    []FieldError `json:"errors,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ContextKey string
