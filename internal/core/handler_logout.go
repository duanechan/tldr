package core

import (
	"net/http"
)

func (t *TLDR) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("REFRESH_TOKEN")
	if err != nil {
		t.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid/missing cookie",
		)
		return
	}

	if err = t.Queries.RevokeRefreshToken(r.Context(), cookie.Value); err != nil {
		t.Logger.Error("Failed to revoke refresh token", "error", err.Error())
		t.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to revoke refresh token",
		)
		return
	}

	t.jsonResponse(w, http.StatusNoContent, nil)
}
