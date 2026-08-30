// Package handlers
package handlers

import (
	"fmt"
	"net/http"
	"quicknotes/internal/models"
	render "quicknotes/internal/render"
	"quicknotes/internal/repositories"
	"strconv"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/csrf"
)

type noteHandler struct {
	render  *render.RenderTenplate
	session *scs.SessionManager
	repo    repositories.NoteRepository
}

func NewNoteHandler(render *render.RenderTenplate, session *scs.SessionManager, repo repositories.NoteRepository) *noteHandler {
	return &noteHandler{render: render, session: session, repo: repo}
}

func (nh *noteHandler) getUserIDFromSession(r *http.Request) int64 {
	return nh.session.GetInt64(r.Context(), "userID")
}

func (nh *noteHandler) NoteList(w http.ResponseWriter, r *http.Request) error {
	notes, err := nh.repo.List(r.Context(), int(nh.getUserIDFromSession(r)))
	if err != nil {
		return err
	}
	return nh.render.RenderPage(w, r, http.StatusOK, "note-home.html", newNoteResponseFromNoteList(notes))
}

func (nh *noteHandler) NoteView(w http.ResponseWriter, r *http.Request) error {
	idParam := r.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return err
	}
	note, err := nh.repo.GetByID(r.Context(), int(nh.getUserIDFromSession(r)), id)
	if err != nil {
		return err
	}
	return nh.render.RenderPage(w, r, http.StatusOK, "note-view.html", newNoteReponseFromNote(note))
}

func (nh *noteHandler) NoteNew(w http.ResponseWriter, r *http.Request) error {
	data := newNoteRequest(nil)
	data.CSRFField = csrf.TemplateField(r)
	return nh.render.RenderPage(w, r, http.StatusOK, "note-new.html", data)
}

func (nh *noteHandler) NoteSave(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
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
			nh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "note-edit.html", data)
		} else {
			nh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "note-new.html", data)
		}
		return nil
	}

	var note *models.Note
	if id > 0 {
		note, err = nh.repo.Update(r.Context(), int(nh.getUserIDFromSession(r)), id, title, content, color)
	} else {
		note, err = nh.repo.Create(r.Context(), int(nh.getUserIDFromSession(r)), title, content, color)
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
	err = nh.repo.Delete(r.Context(), int(nh.getUserIDFromSession(r)), id)
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
	// note, err := nh.repo.GetByID(r.Context(), int(nh.getUserIDFromSession(r), id)))
	note, err := nh.repo.GetByID(r.Context(), int(nh.getUserIDFromSession(r)), id)
	if err != nil {
		return err
	}
	return nh.render.RenderPage(w, r, http.StatusOK, "note-edit.html", newNoteRequest(note))
}
