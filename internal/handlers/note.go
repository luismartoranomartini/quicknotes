// Package handlers
package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"quicknotes/internal/apperror"
	"text/template"
)

type noteHandler struct{}

func NewNoteHandler() *noteHandler {
	return &noteHandler{}
}

func (nh *noteHandler) Notelist(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return errors.New("aconteceu alguma coisa")
	}

	files := []string{
		"views/templates/base.html",
		"views/templates/pages/home.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return errors.New("aconteceu um erro")
	}
	slog.Info("Executou o handler /")
	return tmpl.ExecuteTemplate(w, "base", nil)
}

func (nh *noteHandler) NoteView(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		return apperror.WithStatus(errors.New("anotação é obrigatória"), http.StatusBadRequest)
	}
	if id == "0" {
		return apperror.WithStatus(errors.New("anotação 0 não encontrada"), http.StatusNotFound)
		// return handlers.ErrNotFound
	}
	files := []string{
		"views/templates/base.html",
		"views/templates/pages/note-view.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return errors.New("aconteceu um erro ao executar a página")
	}
	return tmpl.ExecuteTemplate(w, "base", id)
}

func (nh *noteHandler) NoteNew(w http.ResponseWriter, r *http.Request) error {
	files := []string{
		"views/templates/base.html",
		"views/templates/pages/note-new.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return errors.New("aconteceu um erro")
	}
	return tmpl.ExecuteTemplate(w, "base", nil)
}

func (nh *noteHandler) NoteCreate(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		return errors.New("aconteceu um erro")
	}
	fmt.Fprint(w, "Criando uma nova nota")
	return nil
}
