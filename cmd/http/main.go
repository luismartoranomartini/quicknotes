package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"quicknotes/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	config := loadConfig()

	slog.SetDefault(newLogger(os.Stderr, config.GetLevelLog()))

	slog.Info(fmt.Sprintf("Servidor rodando na porta %s", config.ServerPort))

	staticHandler := http.FileServer(http.Dir("views/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// noteHandler := handlers.NewNoteHandler()
	noteHandler := handlers.NewNoteHandler()

	mux.Handle("/", handlers.HandlerWithError(noteHandler.Notelist))
	mux.Handle("/note/view", handlers.HandlerWithError(noteHandler.NoteView))
	mux.Handle("/note/new", handlers.HandlerWithError(noteHandler.NoteNew))
	mux.Handle("/note/create", handlers.HandlerWithError(noteHandler.NoteCreate))

	err := http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), mux)
	if err != nil {
		panic(err)
	}
}
