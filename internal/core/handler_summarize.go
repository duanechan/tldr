package core

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"google.golang.org/genai"
)

func (t *TLDR) SummarizeFile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	file, header, err := r.FormFile("document")
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
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
		errorResponse(w, http.StatusInternalServerError, err.Error())
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
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := t.Client.Models.GenerateContent(
		r.Context(),
		t.Config.APIModel,
		genai.Text(req.Text),
		t.Model,
	)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	duration := time.Since(start)

	jsonResponse(w, http.StatusOK, SummarizeResponse{
		Response: result.Text(),
		Duration: duration.Milliseconds(),
	})
}
