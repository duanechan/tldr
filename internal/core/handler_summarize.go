package core

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/duanechan/tldr/internal/database"
	"github.com/golang-jwt/jwt/v5"
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

func (t *TLDR) SummarizeFile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	claims, ok := r.Context().Value(claimsKey).(*jwt.RegisteredClaims)
	if !ok {
		t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid claims")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to parse user ID")
		return
	}

	if err := r.ParseMultipartForm(maxUploadMemory); err != nil {
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
	if !slices.Contains(allowedFileTypes, mimeType) {
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
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		return
	}

	var response SummarizeResponse
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
	claims, ok := r.Context().Value(claimsKey).(*jwt.RegisteredClaims)
	if !ok {
		t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid claims")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to parse user ID")
		return
	}

	var req SummarizeTextRequest
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
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		return
	}

	var response SummarizeResponse
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

func (t *TLDR) insertTLDR(ctx context.Context, userId uuid.UUID, response SummarizeResponse) error {
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
