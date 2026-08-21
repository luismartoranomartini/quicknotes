// Package handlers
package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"quicknotes/internal/apperror"
	"quicknotes/internal/repositories"
	"strconv"
)

type noteHandler struct {
	repo repositories.NoteRepository
}

func NewNoteHandler(repo repositories.NoteRepository) *noteHandler {
	return &noteHandler{repo: repo}
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
	notes, err := nh.repo.List(r.Context())
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, "base", newNoteResponseFromNoteList(notes))
}

func (nh *noteHandler) NoteView(w http.ResponseWriter, r *http.Request) error {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		return apperror.WithStatus(errors.New("anotação é obrigatória"), http.StatusBadRequest)
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return err
	}
	files := []string{
		"views/templates/base.html",
		"views/templates/pages/note-view.html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return ErrInternal
	}

	note, err := nh.repo.GetByID(r.Context(), id)
	if err != nil {
		return err
	}
	buff := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(buff, "base", newNoteReponseFromNote(note))
	if err != nil {
		return ErrInternal
	}
	buff.WriteTo(w)
	return nil
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
	return tmpl.ExecuteTemplate(w, "base", newNoteRequest())
}

func (nh *noteHandler) NoteCreate(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		return errors.New("aconteceu um erro")
	}
	err := r.ParseForm()
	if err != nil {
		return err
	}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	color := r.PostForm.Get("color")
	// title1 := r.PostFormValue("title")

	note, err := nh.repo.Create(r.Context(), title, content, color)
	if err != nil {
		return err
	}

	http.Redirect(w, r, fmt.Sprintf("/note/view?id=%d", note.ID.Int), http.StatusSeeOther)

	return nil
}

func (nh *noteHandler) NoteDelete(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", http.MethodPost)

		return errors.New("aconteceu um erro")
	}
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		return apperror.WithStatus(errors.New("anotação é obrigatória"), http.StatusBadRequest)
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return err
	}
	err = nh.repo.Delete(r.Context(), id)
	if err != nil {
		return err
	}
	return nil
}
