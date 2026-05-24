package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/db/ping", nil)
	rec := httptest.NewRecorder()
	Ping(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}
