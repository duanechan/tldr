package core

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
	Response string `json:"response"`
	Duration int64  `json:"duration"`
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

type contextKey string

const (
	claimsKey    contextKey = "claims"
	requestIdKey contextKey = "requestId"
)
