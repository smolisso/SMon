package api

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeDB struct{ err error }

func (f fakeDB) Ready(context.Context) error { return f.err }

func TestReadinessDoesNotConfuseLivenessWithDatabaseReadiness(t *testing.T) {
	for _, tc := range []struct {
		name string
		db   Readiness
		code int
		body string
	}{
		{"unconfigured", nil, 503, "not_configured"},
		{"unreachable", fakeDB{errors.New("private connection details")}, 503, "unavailable_or_unmigrated"},
		{"ready", fakeDB{}, 200, `"status":"ready"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := newHandler(tc.db)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, httptest.NewRequest("GET", "/api/v1/ready", nil))
			if w.Code != tc.code || !strings.Contains(w.Body.String(), tc.body) || strings.Contains(w.Body.String(), "private") {
				t.Fatalf("unexpected readiness: %d %s", w.Code, w.Body.String())
			}
			w = httptest.NewRecorder()
			h.ServeHTTP(w, httptest.NewRequest("GET", "/api/v1/health", nil))
			if w.Code != 200 {
				t.Fatalf("health failed: %d", w.Code)
			}
		})
	}
}

func TestEmbeddedUIAndUnknownAPI(t *testing.T) {
	h := newHandler(nil)
	for _, tc := range []struct {
		path    string
		code    int
		content string
	}{
		{"/", 200, `<div id="root">`},
		{"/api/v1/missing", 404, "endpoint_not_found"},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
		if w.Code != tc.code || !strings.Contains(w.Body.String(), tc.content) {
			t.Fatalf("%s: %d %s", tc.path, w.Code, w.Body.String())
		}
	}
}
