package core

type Page[T any] struct {
	Next    PageCursor `json:"next,omitempty"`
	Results []T        `json:"results"`
}

type PageCursor string
type PageLimit int64

const (
	PageCursorSeparator = "|"
	DefaultPageCursor   = "MjUzNDAyMjE0NDAwfGZmZmZmZmZmLWZmZmYtZmZmZi1mZmZmLWZmZmZmZmZmZmZmZg"
	DefaultPageLimit    = "10"
	MaxPageLimit        = 100
)

type HealthResponse struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	AccessToken string `json:"access_token"`
}

type registerRequest struct {
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

const sqliteUniqueConstraint = 2067

const (
	requestIdKey = "requestId"
)
