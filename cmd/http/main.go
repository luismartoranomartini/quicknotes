package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"quicknotes/internal/handlers"
	"quicknotes/internal/repositories"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config := loadConfig()
	slog.SetDefault(newLogger(os.Stderr, config.GetLevelLog()))

	dbpool, err := pgxpool.New(context.Background(), config.DBConnURL)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("Conexão com o banco aconteceu com sucesso")
	defer dbpool.Close()

	mux := http.NewServeMux()

	sesseionManager := scs.New()
	sesseionManager.Lifetime = time.Hour

	slog.Info(fmt.Sprintf("Servidor rodando na porta %s", config.ServerPort))

	staticHandler := http.FileServer(http.Dir("views/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", staticHandler))

	noteRepo := repositories.NewNoteRepository(dbpool)
	userRepo := repositories.NewUserRepository(dbpool)

	noteHandler := handlers.NewNoteHandler(noteRepo)
	userHandler := handlers.NewUserHandler(sesseionManager, userRepo)

	mux.Handle("/", handlers.HandlerWithError(noteHandler.NoteList))
	mux.Handle("GET /note/{id}", handlers.HandlerWithError(noteHandler.NoteView))
	mux.Handle("GET /note/new", handlers.HandlerWithError(noteHandler.NoteNew))
	mux.Handle("POST /note", handlers.HandlerWithError(noteHandler.NoteSave))
	mux.Handle("DELETE /note/{id}", handlers.HandlerWithError(noteHandler.NoteDelete))
	mux.Handle("GET /note/{id}/edit", handlers.HandlerWithError(noteHandler.NoteEdit))

	mux.Handle("GET /user/signup", handlers.HandlerWithError(userHandler.SignupForm))
	mux.Handle("POST /user/signup", handlers.HandlerWithError(userHandler.Signup))

	mux.Handle("GET /user/signin", handlers.HandlerWithError(userHandler.SigninForm))
	mux.Handle("POST /user/signin", handlers.HandlerWithError(userHandler.Signin))

	mux.Handle("GET /me", handlers.HandlerWithError(userHandler.Me))

	// middleware
	mux.Handle("GET /confirmation/{token}", handlers.HandlerWithError(userHandler.Confirm))

	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), sesseionManager.LoadAndSave(csrfMiddleware(mux))); err != nil {
		panic(err)
	}
}
