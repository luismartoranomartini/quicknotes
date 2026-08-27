// Package handlers
package handlers

import (
	"fmt"
	"net/http"
	"quicknotes/internal/models"
	"quicknotes/internal/repositories"
	"strconv"
	"strings"

	"github.com/gorilla/csrf"
)

type noteHandler struct {
	repo repositories.NoteRepository
}

func NewNoteHandler(repo repositories.NoteRepository) *noteHandler {
	return &noteHandler{repo: repo}
}

func (nh *noteHandler) NoteList(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return ErrNotFound
	}

	notes, err := nh.repo.List(r.Context())
	if err != nil {
		return err
	}
	return render(w, http.StatusOK, "home.html", newNoteResponseFromNoteList(notes))
}

func (nh *noteHandler) NoteView(w http.ResponseWriter, r *http.Request) error {
	idParam := r.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return err
	}
	note, err := nh.repo.GetByID(r.Context(), id)
	if err != nil {
		return err
	}
	return render(w, http.StatusOK, "note-view.html", newNoteReponseFromNote(note))
}

func (nh *noteHandler) NoteNew(w http.ResponseWriter, r *http.Request) error {
	data := newNoteRequest(nil)
	data.CSRFField = csrf.TemplateField(r)
	return render(w, http.StatusOK, "note-new.html", data)
}

func (nh *noteHandler) NoteSave(w http.ResponseWriter, r *http.Request) error {
	// _, err := r.Cookie("session")
	// if err != nil {
	// 	http.Redirect(w, r, "/user/signin", http.StatusTemporaryRedirect)
	// }

	err := r.ParseForm()
	if err != nil {
		return err
	}
	fmt.Println(r.PostForm.Get("gorilla.csrf.Token"))
	idParam := r.PostForm.Get("id")
	id, _ := strconv.Atoi(idParam)
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	color := r.PostForm.Get("color")
	// title1 := r.PostFormValue("title") -> outra forma de fazer

	data := newNoteRequest(nil)
	data.ID = id
	data.Title = title
	data.Content = content
	data.Color = color

	if strings.TrimSpace(content) == "" {
		data.AddFieldError("content", "Conteúdo é obrigatório")
	}

	if !data.Valid() {
		if id > 0 {
			render(w, http.StatusUnprocessableEntity, "note-edit.html", data)
		} else {
			render(w, http.StatusUnprocessableEntity, "note-new.html", data)
		}
		return nil
	}

	var note *models.Note
	if id > 0 {
		note, err = nh.repo.Update(r.Context(), id, title, content, color)
	} else {
		note, err = nh.repo.Create(r.Context(), title, content, color)
	}

	if err != nil {
		return err
	}
	http.Redirect(w, r, fmt.Sprintf("/note/%d", note.ID.Int), http.StatusSeeOther)
	return nil
}

func (nh *noteHandler) NoteDelete(w http.ResponseWriter, r *http.Request) error {
	idParam := r.PathValue("id")
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

func (nh *noteHandler) NoteEdit(w http.ResponseWriter, r *http.Request) error {
	idParam := r.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return err
	}
	note, err := nh.repo.GetByID(r.Context(), id)
	if err != nil {
		return err
	}
	return render(w, http.StatusOK, "note-edit.html", newNoteRequest(note))
}
