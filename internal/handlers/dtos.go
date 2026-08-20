package handlers

import "quicknotes/internal/models"

type NoteResponse struct {
	ID      int
	Title   string
	Content string
}

func newNoteReponseFromNote(note *models.Note) (res NoteResponse) {
	res.ID = int(note.ID.Int.Int64())
	res.Title = note.Title.String
	res.Content = note.Content.String

	return
}

func newNoteResponseFromNoteList(notes []models.Note) (res []NoteResponse) {
	for _, note := range notes {
		res = append(res, newNoteReponseFromNote(&note))
	}
	return
}
