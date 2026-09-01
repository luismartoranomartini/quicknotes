package main

import (
	"net/http"
	"quicknotes/internal/handlers"
	"quicknotes/internal/mailer"
	render "quicknotes/internal/render"
	"quicknotes/internal/repositories"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func LoadRoutes(sesseionManager *scs.SessionManager, mail mailer.MailService, dbpool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	staticHandler := http.FileServer(http.Dir("views/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", staticHandler))

	noteRepo := repositories.NewNoteRepository(dbpool)
	userRepo := repositories.NewUserRepository(dbpool)

	render := render.NewRender(sesseionManager)

	noteHandler := handlers.NewNoteHandler(render, sesseionManager, noteRepo)
	userHandler := handlers.NewUserHandler(render, sesseionManager, mail, userRepo)
	authMiddleware := handlers.NewAuthMiddleware(sesseionManager)

	mux.HandleFunc("GET /", handlers.NewHomeHandler(render).HomeHandler)

	mux.Handle("GET /note", authMiddleware.RequireAuth(handlers.HandlerWithError(noteHandler.NoteList)))
	mux.Handle("GET /note/{id}", handlers.HandlerWithError(noteHandler.NoteView))
	mux.Handle("GET /note/new", authMiddleware.RequireAuth(handlers.HandlerWithError(noteHandler.NoteNew)))
	mux.Handle("POST /note", authMiddleware.RequireAuth(handlers.HandlerWithError(noteHandler.NoteSave)))
	mux.Handle("DELETE /note/{id}", authMiddleware.RequireAuth(handlers.HandlerWithError(noteHandler.NoteDelete)))
	mux.Handle("GET /note/{id}/edit", authMiddleware.RequireAuth(handlers.HandlerWithError(noteHandler.NoteEdit)))

	mux.Handle("GET /user/signup", handlers.HandlerWithError(userHandler.SignupForm))
	mux.Handle("POST /user/signup", handlers.HandlerWithError(userHandler.Signup))

	mux.Handle("GET /user/signin", handlers.HandlerWithError(userHandler.SigninForm))
	mux.Handle("POST /user/signin", handlers.HandlerWithError(userHandler.Signin))

	mux.Handle("GET /user/signout", handlers.HandlerWithError(userHandler.Signout))

	mux.Handle("GET /me", authMiddleware.RequireAuth(handlers.HandlerWithError(userHandler.Me)))

	// middleware
	mux.Handle("GET /confirmation/{token}", handlers.HandlerWithError(userHandler.Confirm))

	return mux
}
