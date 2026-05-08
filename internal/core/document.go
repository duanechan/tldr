package core

import (
	"io"
	"net/http"
)

type DocumentData struct {
	Filename  string
	FileBytes []byte
	MIMEType  string
}

func GetDocumentData(w http.ResponseWriter, r *http.Request) (*DocumentData, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, err
	}

	file, header, err := r.FormFile("document")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &DocumentData{
		Filename:  header.Filename,
		FileBytes: fileBytes,
		MIMEType:  mimeType,
	}, nil
}
