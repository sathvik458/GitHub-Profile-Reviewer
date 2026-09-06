package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthReportsOK(t *testing.T) {
	rec := httptest.NewRecorder()

	Health(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var response healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode health response: %v", err)
	}

	if response.Status != "ok" {
		t.Fatalf("status = %q, want %q", response.Status, "ok")
	}
}

func TestHealthRejectsInvalidMethod(t *testing.T) {
	rec := httptest.NewRecorder()

	Health(rec, httptest.NewRequest(http.MethodPost, "/healthz", nil))

	assertErrorResponse(t, rec, http.StatusMethodNotAllowed, "method not allowed")

	if allow := rec.Header().Get("Allow"); allow != http.MethodGet {
		t.Fatalf("Allow header = %q, want %q", allow, http.MethodGet)
	}
}
