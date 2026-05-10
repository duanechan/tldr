package core

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"google.golang.org/genai"
)

type HealthResponse struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}

func (t *TLDR) Health(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, HealthResponse{
		Status: "OK",
		Uptime: time.Since(t.startedAt).Round(time.Second).String(),
	})
}

type SummarizeResponse struct {
	Response string `json:"response"`
	Filename string `json:"filename"`
	FileType string `json:"file_type"`
	Duration int64  `json:"duration"`
}

func (t *TLDR) SummarizeDocument(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if err := r.ParseMultipartForm(10 << 2); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	file, header, err := r.FormFile("document")
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	filename := header.Filename
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
		Filename: filename,
		FileType: mimeType,
		Duration: duration.Milliseconds(),
	})
}

type SummarizeTextRequest struct {
	Text string `json:"text"`
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
		Filename: "",
		FileType: "text/plain",
		Duration: duration.Milliseconds(),
	})
}
