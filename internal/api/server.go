package api

import (
	"context"
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/smolisso/SMon/internal/postgres"
	"github.com/smolisso/SMon/web"
)

type Readiness interface{ Ready(context.Context) error }

func New(store *postgres.Store) http.Handler {
	var readiness Readiness
	if store != nil {
		readiness = store
	}
	return newHandler(readiness)
}

func newHandler(db Readiness) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusOK, map[string]string{"status": "ok", "name": "SMon", "stage": "foundation"})
	})
	mux.HandleFunc("GET /api/v1/ready", func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			respond(w, http.StatusServiceUnavailable, map[string]string{"status": "not_ready", "database": "not_configured"})
			return
		}
		if err := db.Ready(r.Context()); err != nil {
			respond(w, http.StatusServiceUnavailable, map[string]string{"status": "not_ready", "database": "unavailable_or_unmigrated"})
			return
		}
		respond(w, http.StatusOK, map[string]string{"status": "ready", "database": "ready"})
	})
	mux.HandleFunc("GET /api/", func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusNotFound, map[string]string{"error": "endpoint_not_found"})
	})
	assets, err := fs.Sub(web.Files, "dist")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /", http.FileServer(http.FS(assets)))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'")
		mux.ServeHTTP(w, r)
	})
}

func respond(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(value)
}
