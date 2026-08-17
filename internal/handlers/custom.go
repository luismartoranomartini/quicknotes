package handlers

import (
	"errors"
	"net/http"
	"quicknotes/internal/apperror"
	"text/template"
)

var ErrNotFound = apperror.WithStatus(errors.New("not found"), http.StatusNotFound)

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

func (f HandlerWithError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := f(w, r); err != nil {
		var statusErr apperror.StatusError
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
