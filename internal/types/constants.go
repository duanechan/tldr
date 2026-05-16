package types

import "time"

const (
	Prompt = `
	You are a document summarizer.
	
	Summarize the documents, text files, raw text, or images that are provided
	to you.
	
	Extract its keypoints, if necessary, and ensure the user understands what
	the contents are.
	
	Keep the content brief and simple (250-500 words).
	
	Don't make the title too verbose.
	
	Do not ask for more input.
	
	No markdown.
	
	ONLY RETURN a response in this JSON format:
	{
		"title": string,
		"content": string
	}
	`

	MaxUploadMemory int64 = 10 << 20

	DefaultPageLimit  = "10"
	DefaultPageCursor = "9999-12-31T00:00:00Z"

	ClaimsKey    ContextKey = "claims"
	RequestIdKey ContextKey = "requestId"

	MinimumUsernameLength = 3
	MinimumPasswordLength = 8

	RateLimitInterval = 1 * time.Minute
	BurstLimit        = 5

	UniqueConstraintCode = 2067
)

var AllowedFileTypes = []string{
	"application/pdf",
	"image/png",
	"image/jpeg",
	"image/gif",
	"text/plain",
}
