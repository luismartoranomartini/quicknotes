package handlers

import (
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type AuthMiddleware struct {
	session *scs.SessionManager
}

func NewAuthMiddleware(session *scs.SessionManager) *AuthMiddleware {
	return &AuthMiddleware{
		session: session,
	}
}

func (ah *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := ah.session.GetInt64(r.Context(), "userID")
		if userID == 0 {
			slog.Warn("usuário não está logado")
			http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
