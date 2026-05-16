package core

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/duanechan/tldr/internal/auth"
	"github.com/duanechan/tldr/internal/database"
	"github.com/duanechan/tldr/internal/types"
	"github.com/google/uuid"
	"google.golang.org/genai"
)

func (t *TLDR) SummarizeFile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	userId, err := auth.GetUserID(r.Context())
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid claims")
		return
	}

	if err := r.ParseMultipartForm(types.MaxUploadMemory); err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Invalid or missing multipart form data")
		return
	}

	file, header, err := r.FormFile("document")
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Invalid or missing document field")
		return
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")
	if !slices.Contains(types.AllowedFileTypes, mimeType) {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "File type not supported")
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		return
	}

	parts := []*genai.Part{
		{InlineData: &genai.Blob{MIMEType: mimeType, Data: fileBytes}},
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := t.Client.Models.GenerateContent(
		r.Context(),
		t.Config.APIModel,
		contents,
		t.Model,
	)
	if err, exists := errors.AsType[genai.APIError](err); exists {
		switch err.Code {
		case http.StatusTooManyRequests:
			t.errorResponse(w, r.Context(), http.StatusTooManyRequests, "Too many requests, try again later")
		case http.StatusRequestEntityTooLarge:
			t.errorResponse(w, r.Context(), http.StatusRequestEntityTooLarge, "File is too large")
		default:
			t.Logger.Error("Failed to summarize text", "error", err.Error())
			t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		}
		return
	}

	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		return
	}

	var response types.SummarizeResponse
	if err := json.Unmarshal([]byte(result.Text()), &response); err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		return
	}

	if err = t.insertTLDR(r.Context(), userId, response); err != nil {
		t.Logger.Error("Failed to create TLDR", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to create TLDR")
		return
	}

	response.Duration = time.Since(start).Milliseconds()
	t.jsonResponse(w, http.StatusOK, response)
}

func (t *TLDR) SummarizeText(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	userId, err := auth.GetUserID(r.Context())
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid claims")
		return
	}

	var req types.SummarizeTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Invalid request body")
		return
	}

	cleanedText := strings.TrimSpace(req.Text)
	if cleanedText == "" {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Text is required")
		return
	}

	result, err := t.Client.Models.GenerateContent(
		r.Context(),
		t.Config.APIModel,
		genai.Text(cleanedText),
		t.Model,
	)
	if err, exists := errors.AsType[genai.APIError](err); exists {
		switch err.Code {
		case http.StatusTooManyRequests:
			t.errorResponse(w, r.Context(), http.StatusTooManyRequests, "Too many requests, try again later")
		case http.StatusRequestEntityTooLarge:
			t.errorResponse(w, r.Context(), http.StatusRequestEntityTooLarge, "File is too large")
		default:
			t.Logger.Error("Failed to summarize text", "error", err.Error())
			t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		}
		return
	}

	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		return
	}

	var response types.SummarizeResponse
	if err = json.Unmarshal([]byte(result.Text()), &response); err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		return
	}

	if err = t.insertTLDR(r.Context(), userId, response); err != nil {
		t.Logger.Error("Failed to create TLDR", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to create TLDR")
		return
	}

	response.Duration = time.Since(start).Milliseconds()
	t.jsonResponse(w, http.StatusOK, response)
}

func (t *TLDR) insertTLDR(ctx context.Context, userId uuid.UUID, response types.SummarizeResponse) error {
	tldrId, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	_, err = t.Queries.CreateTLDR(ctx, database.CreateTLDRParams{
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
