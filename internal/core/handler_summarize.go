package core

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

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
	if err := r.ParseMultipartForm(maxUploadMemory); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid or missing multipart form data")
		return
	}

	file, header, err := r.FormFile("document")
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid or missing document field")
		return
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")
	if !slices.Contains(allowedFileTypes, mimeType) {
		errorResponse(w, http.StatusBadRequest, "File type not supported")
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong")
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
		errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	duration := time.Since(start)

	jsonResponse(w, http.StatusOK, SummarizeResponse{
		Response: result.Text(),
		Duration: duration.Milliseconds(),
	})
}

func (t *TLDR) SummarizeText(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req SummarizeTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cleanedText := strings.TrimSpace(req.Text)
	if cleanedText == "" {
		errorResponse(w, http.StatusBadRequest, "Text is required")
		return
	}

	result, err := t.Client.Models.GenerateContent(
		r.Context(),
		t.Config.APIModel,
		genai.Text(cleanedText),
		t.Model,
	)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	duration := time.Since(start)

	jsonResponse(w, http.StatusOK, SummarizeResponse{
		Response: result.Text(),
		Duration: duration.Milliseconds(),
	})
}
