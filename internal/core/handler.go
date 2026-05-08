package core

import (
	"net/http"
	"time"

	"google.golang.org/genai"
)

func (t *TLDR) Health(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 200, struct {
		Status string
	}{Status: "OK"})
}

type SummarizeResponse struct {
	Response string `json:"response"`
	Filename string `json:"filename"`
	FileType string `json:"file_type"`
	Duration int64  `json:"duration"`
}

func (t *TLDR) SummarizeDocument(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	data, err := GetDocumentData(w, r)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	parts := []*genai.Part{
		{InlineData: &genai.Blob{MIMEType: data.MIMEType, Data: data.FileBytes}},
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
		Filename: data.Filename,
		FileType: data.MIMEType,
		Duration: duration.Milliseconds(),
	})
}
