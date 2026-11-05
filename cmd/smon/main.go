package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smolisso/SMon/internal/api"
	"github.com/smolisso/SMon/internal/postgres"
)

func main() {
	if err := run(); err != nil {
		slog.Error("SMon stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	command := "serve"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if len(os.Args) > 2 || (command != "serve" && command != "migrate") {
		return errors.New("usage: smon [serve|migrate]")
	}
	db, err := postgres.Open(ctx, os.Getenv("SMON_DATABASE_URL"))
	if err != nil {
		return err
	}
	if db != nil {
		defer db.Close()
	}
	if command == "migrate" {
		if db == nil {
			return errors.New("SMON_DATABASE_URL is required for migrations")
		}
		if err := db.Migrate(ctx); err != nil {
			return err
		}
		slog.Info("database migrations applied")
		return nil
	}
	addr := os.Getenv("SMON_LISTEN")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	server := &http.Server{Addr: addr, Handler: api.New(db), ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	result := make(chan error, 1)
	go func() { result <- server.ListenAndServe() }()
	slog.Info("SMon starting", "listen", addr, "database_configured", db != nil)
	select {
	case err := <-result:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdown)
	}
}
