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
