package core

import (
	"encoding/json"
	"log"
	"net/http"
)

func errorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(struct {
		Error string
		Code  int
	}{
		Error: message,
		Code:  code,
	})
	if err != nil {
		log.Printf("marshal error: %v", err)
		code = http.StatusInternalServerError
		body = []byte(`{"error":"Something went wrong"}`)
	}

	w.WriteHeader(code)

	if n, err := w.Write(body); err != nil {
		log.Printf("write error: %v", err)
	} else if n < len(body) {
		log.Printf("short write: wrote %d of %d bytes", n, len(body))
	}
}

func jsonResponse(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("marshal error: %v", err)
		code = http.StatusInternalServerError
		body = []byte(`{"error":"Something went wrong"}`)
	}

	w.WriteHeader(code)

	if n, err := w.Write(body); err != nil {
		log.Printf("write error: %v", err)
	} else if n < len(body) {
		log.Printf("short write: wrote %d of %d bytes", n, len(body))
	}
}
