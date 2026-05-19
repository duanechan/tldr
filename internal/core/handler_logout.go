package core

import (
	"net/http"
	"time"
)

// Logout handles logout requests by revoking the refresh token attached
// in the cookie, and clears it.
func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("REFRESH_TOKEN")
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid/missing cookie",
		)
		return
	}

	if err = a.Queries.RevokeRefreshToken(r.Context(), cookie.Value); err != nil {
		a.Logger.Error("Failed to revoke refresh token", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to revoke refresh token",
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "REFRESH_TOKEN",
		Value:   "",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
	})

	a.jsonResponse(w, http.StatusNoContent, nil)
}
