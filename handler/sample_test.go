package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{"happy path", "/?name=akira", http.StatusOK, "Hello akira"},
		{"missing name", "/", http.StatusBadRequest, "set name parameter."},
		{"not found path", "/foo", http.StatusNotFound, "not found."},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()
			Index(rec, req)
			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if !strings.Contains(rec.Body.String(), tc.wantBody) {
				t.Errorf("body = %q, want substring %q", rec.Body.String(), tc.wantBody)
			}
		})
	}
}
