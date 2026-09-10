package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"quicknotes/internal/apperror"
	"quicknotes/internal/render"
	"quicknotes/internal/repositories"

	"github.com/alexedwards/scs/v2"
)

type AuthMiddleware struct {
	session *scs.SessionManager
}

var ErrNotFound = apperror.WithStatus(errors.New("not found"), http.StatusNotFound)
var ErrInternal = apperror.WithStatus(errors.New("aconteceu um erro ao executar a página"), http.StatusInternalServerError)

func NewAuthMiddleware(session *scs.SessionManager) *AuthMiddleware {
	return &AuthMiddleware{
		session: session,
	}
}

type errorHandlerMiddleware struct {
	render *render.RenderTemplate
}

// função de instanciação
func NewErrorHandleMiddleware(render *render.RenderTemplate) *errorHandlerMiddleware {
	return &errorHandlerMiddleware{render: render}
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

func (em *errorHandlerMiddleware) HandleError(next func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			var statusErr apperror.StatusError
			var repoErr repositories.RepositoryError
			if errors.As(err, &statusErr) {
				if statusErr.StatusCode() == http.StatusNotFound {
					em.render.RenderPage(w, r, http.StatusNotFound, "404.html", nil)
					return
				}
			}
			if errors.As(err, &repoErr) {
				slog.Error(err.Error())
				em.render.RenderPage(w, r, http.StatusInternalServerError, "generic-error.html", "aconteceu um erro ao executar uma operação")
				return
			}
			slog.Error(err.Error())
			em.render.RenderPage(w, r, http.StatusInternalServerError, "generic-error.html", err.Error())
		}
	})
}
