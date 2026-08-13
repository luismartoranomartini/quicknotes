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

	slog.Info(fmt.Sprintf("DBPASSWORD = %s", config.DBPassword))

	slog.Info(fmt.Sprintf("Servidor rodando na porta %s", config.ServerPort))

	staticHandler := http.FileServer(http.Dir("views/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// noteHandler := handlers.NewNoteHandler()
	noteHandler := handlers.NewNoteHandler()

	mux.HandleFunc("/", noteHandler.Notelist)
	mux.HandleFunc("/note/view", noteHandler.NoteView)
	mux.HandleFunc("/note/new", noteHandler.NoteNew)
	mux.HandleFunc("/note/create", noteHandler.NoteCreate)

	http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), mux)
}
