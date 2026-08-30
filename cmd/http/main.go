package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
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

	sesseionManager := scs.New()
	sesseionManager.Lifetime = time.Hour
	sesseionManager.Store = pgxstore.New(dbpool)
	//Limpa as sessões expiradas da tabbela
	pgxstore.NewWithCleanupInterval(dbpool, 30*time.Minute)

	// slog.Info(fmt.Sprintf("Servidor rodando na porta %s", config.ServerPort))

	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"))
	mux := LoadRoutes(sesseionManager, dbpool)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), sesseionManager.LoadAndSave(csrfMiddleware(mux))); err != nil {
		panic(err)
	}
}
