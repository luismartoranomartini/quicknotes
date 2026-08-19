package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"quicknotes/internal/handlers"
	"quicknotes/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config := loadConfig()
	mux := http.NewServeMux()

	slog.SetDefault(newLogger(os.Stderr, config.GetLevelLog()))

	dbpool, err := pgxpool.New(context.Background(), config.DBConnUrl)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("Conexão com o banco aconteceu com sucesso")
	defer dbpool.Close()

	slog.Info(fmt.Sprintf("Servidor rodando na porta %s", config.ServerPort))

	staticHandler := http.FileServer(http.Dir("views/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	noteRepo := repositories.NewNoteRepository(dbpool)

	err = noteRepo.Delete(3)

	if err != nil {
		slog.Error(err.Error())
	}
	fmt.Println("nota deletada com sucesso")

	noteHandler := handlers.NewNoteHandler()

	mux.Handle("/", handlers.HandlerWithError(noteHandler.Notelist))
	mux.Handle("/note/view", handlers.HandlerWithError(noteHandler.NoteView))
	mux.Handle("/note/new", handlers.HandlerWithError(noteHandler.NoteNew))
	mux.Handle("/note/create", handlers.HandlerWithError(noteHandler.NoteCreate))

	err = http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), mux)
	if err != nil {
		panic(err)
	}
}
