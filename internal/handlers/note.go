// Package handlers
package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"text/template"
)

type noteHandler struct{}

// NewNoteHandler() função que aplica os handlers
func NewNoteHandler() *noteHandler {
	return &noteHandler{}
}

func (nh *noteHandler) Notelist(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"views/templates/base.html",
		"views/templates/pages/home.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "aconteceu um erro", http.StatusInternalServerError)
		return
	}
	slog.Info("Executou o handler /")
	tmpl.ExecuteTemplate(w, "base", nil)
}

func (nh *noteHandler) NoteView(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "nota nao encotrada", http.StatusNotFound)
		return
	}
	files := []string{
		"views/templates/base.html",
		"views/templates/pages/Note-view.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "aconteceu um erro", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", id)
}

func (nh *noteHandler) NoteNew(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"views/templates/base.html",
		"views/templates/pages/Note-new.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "aconteceu um erro", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", nil)
}

func (nh *noteHandler) NoteCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, "Criando uma nova nota")
}
