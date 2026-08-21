package handlers

import (
	"fmt"
	"quicknotes/internal/models"
)

type NoteResponse struct {
	ID      int
	Title   string
	Content string
	Color   string
}

type NoteRequest struct {
	ID      int
	Title   string
	Content string
	Color   string
	Colors  []string
}

func newNoteRequest() (req NoteRequest) {
	req.Color = "color5"
	for i := 1; i <= 9; i++ {
		req.Colors = append(req.Colors, fmt.Sprintf("color%d", i))
	}
	return
}

func newNoteReponseFromNote(note *models.Note) (res NoteResponse) {
	res.ID = int(note.ID.Int.Int64())
	res.Title = note.Title.String
	res.Content = note.Content.String
	res.Color = note.Color.String

	return
}

func newNoteResponseFromNoteList(notes []models.Note) (res []NoteResponse) {
	for _, note := range notes {
		res = append(res, newNoteReponseFromNote(&note))
	}
	return
}
