package handlers

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"quicknotes/internal/apperror"
	"quicknotes/internal/repositories"
)

var ErrNotFound = apperror.WithStatus(errors.New("not found"), http.StatusNotFound)
var ErrInternal = apperror.WithStatus(errors.New("aconteceu um erro ao executar a página"), http.StatusInternalServerError)

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

func (f HandlerWithError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := f(w, r); err != nil {
		var statusErr apperror.StatusError
		var repoErr repositories.RepositoryError
		if errors.As(err, &statusErr) {
			if statusErr.StatusCode() == http.StatusNotFound {
				files := []string{
					"views/templates/base.html",
					"views/templates/pages/404.html",
				}
				tmpl, err := template.ParseFiles(files...)
				if err != nil {
					http.Error(w, err.Error(), statusErr.StatusCode())
				}
				tmpl.ExecuteTemplate(w, "base", statusErr.Error())
				return

			}

		}
		if errors.As(err, &repoErr) {
			slog.Error(err.Error())
			http.Error(w, "aconteceu um erro ao executar uma operação", http.StatusInternalServerError)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
}
