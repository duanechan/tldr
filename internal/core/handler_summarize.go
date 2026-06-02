package core

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/duanechan/tldr/internal/auth"
	"github.com/duanechan/tldr/internal/database"
	"github.com/duanechan/tldr/internal/validate"
	"github.com/google/uuid"
	"google.golang.org/genai"
)

var allowedFileTypes []string = []string{
	"application/pdf",
	"image/png",
	"image/jpeg",
	"image/gif",
	"text/plain",
}

const (
	maxUploadMemory int64 = 10 << 20
)

type SummarizeResponse struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Flag     string `json:"flag"`
	Duration int64  `json:"duration,omitempty"`
}

// SummarizeFile returns a TLDR of a file.
func (a *App) SummarizeFile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	userId, err := auth.GetUserID(r.Context())
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusUnauthorized,
			"Invalid claims",
		)
		return
	}

	if err := r.ParseMultipartForm(maxUploadMemory); err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid or missing multipart form data",
		)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid or missing file field",
		)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Something went wrong",
		)
		return
	}

	mimeType := http.DetectContentType(fileBytes)
	if !slices.Contains(allowedFileTypes, mimeType) {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"File type not supported",
		)
		return
	}

	parts := []*genai.Part{
		{InlineData: &genai.Blob{MIMEType: mimeType, Data: fileBytes}},
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := a.AI.GenerateContent(
		r.Context(),
		a.Config.APIModel,
		contents,
		ContentConfig,
	)
	if err, exists := errors.AsType[genai.APIError](err); exists {
		switch err.Code {
		case http.StatusTooManyRequests:
			a.errorResponse(
				w,
				r.Context(),
				http.StatusTooManyRequests,
				"Too many requests, try again later",
			)
		case http.StatusRequestEntityTooLarge:
			a.errorResponse(
				w,
				r.Context(),
				http.StatusRequestEntityTooLarge,
				"File is too large",
			)
		default:
			a.Logger.Error("Failed to summarize text", "error", err.Error())
			a.errorResponse(
				w,
				r.Context(),
				http.StatusInternalServerError,
				"Something went wrong",
			)
		}
		return
	}

	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Something went wrong",
		)
		return
	}

	var response SummarizeResponse
	if err := json.Unmarshal([]byte(result.Text()), &response); err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Something went wrong",
		)
		return
	}

	if err = a.insertTLDR(r.Context(), userId, response); err != nil {
		a.Logger.Error("Failed to create TLDR", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to create TLDR",
		)
		return
	}

	response.Duration = time.Since(start).Milliseconds()
	a.jsonResponse(w, http.StatusOK, response)
}

// SummarizeText returns a TLDR of a text.
func (a *App) SummarizeText(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	userId, err := auth.GetUserID(r.Context())
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusUnauthorized,
			"Invalid claims",
		)
		return
	}

	var text string
	if err := json.NewDecoder(r.Body).Decode(&text); err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	text, errs := validate.String(text, validate.NotEmpty())
	if errs != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Text is required",
		)
		return
	}

	result, err := a.AI.GenerateContent(
		r.Context(),
		a.Config.APIModel,
		genai.Text(text),
		ContentConfig,
	)
	if err, exists := errors.AsType[genai.APIError](err); exists {
		switch err.Code {
		case http.StatusTooManyRequests:
			a.errorResponse(
				w,
				r.Context(),
				http.StatusTooManyRequests,
				"Too many requests, try again later",
			)
		case http.StatusRequestEntityTooLarge:
			a.errorResponse(
				w,
				r.Context(),
				http.StatusRequestEntityTooLarge,
				"File is too large",
			)
		default:
			a.Logger.Error("Failed to summarize text", "error", err.Error())
			a.errorResponse(
				w,
				r.Context(),
				http.StatusInternalServerError,
				"Something went wrong",
			)
		}
		return
	}

	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Something went wrong",
		)
		return
	}

	var response SummarizeResponse
	if err = json.Unmarshal([]byte(result.Text()), &response); err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Something went wrong",
		)
		return
	}

	if err = a.insertTLDR(r.Context(), userId, response); err != nil {
		a.Logger.Error("Failed to create TLDR", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to create TLDR",
		)
		return
	}

	response.Duration = time.Since(start).Milliseconds()
	a.jsonResponse(w, http.StatusOK, response)
}

// insertTLDR generates a random UUID for the TLDR, and inserts it based on
// the response.
func (a *App) insertTLDR(
	ctx context.Context,
	userId uuid.UUID,
	response SummarizeResponse,
) error {
	tldrId, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	_, err = a.Queries.CreateTLDR(ctx, database.CreateTLDRParams{
		ID:      tldrId,
		Title:   response.Title,
		Content: response.Content,
		UserID:  userId,
	})
	if err != nil {
		return err
	}

	return nil
}
