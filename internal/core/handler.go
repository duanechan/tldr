package core

import "net/http"

func (t *TLDR) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jsonResponse(w, 200, struct {
		Status string
	}{Status: "OK"})
}
