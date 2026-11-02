//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestErrorResponses(t *testing.T) {
	resetDB(t)

	cases := []struct {
		name       string
		method     string
		path       string
		body       any
		wantStatus int
		wantErr    string
	}{
		{
			name:       "unknown product",
			method:     http.MethodGet,
			path:       "/products/" + uuid.NewString(),
			wantStatus: http.StatusNotFound,
			wantErr:    "product not found",
		},
		{
			name:       "unknown supplier",
			method:     http.MethodGet,
			path:       "/suppliers/" + uuid.NewString(),
			wantStatus: http.StatusNotFound,
			wantErr:    "supplier not found",
		},
		{
			name:       "malformed product id",
			method:     http.MethodGet,
			path:       "/products/not-a-uuid",
			wantStatus: http.StatusBadRequest,
			wantErr:    "invalid UUID",
		},
		{
			name:       "broken json",
			method:     http.MethodPost,
			path:       "/clients",
			body:       rawBody("{"),
			wantStatus: http.StatusBadRequest,
			wantErr:    "invalid JSON format",
		},
		{
			name:   "invalid gender",
			method: http.MethodPost,
			path:   "/clients",
			body: map[string]any{
				"name":      uniqueName("client"),
				"surname":   "Testov",
				"birthdate": "1990-05-10",
				"gender":    "other",
				"address":   addressBody("Russia", "Moscow", "Error Street"),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    "invalid JSON format",
		},
		{
			name:       "underage client",
			method:     http.MethodPost,
			path:       "/clients",
			body:       clientBody(uniqueName("client"), "2020-01-01", addressBody("Russia", "Moscow", "Error Street")),
			wantStatus: http.StatusBadRequest,
			wantErr:    "at least 18 years old",
		},
		{
			name:       "invalid query limit",
			method:     http.MethodGet,
			path:       "/clients?limit=101",
			wantStatus: http.StatusBadRequest,
			wantErr:    "invalid JSON format",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := doJSON(t, tc.method, tc.path, tc.body)
			if status != tc.wantStatus {
				t.Fatalf("status = %d, want %d, body = %s", status, tc.wantStatus, body)
			}
			if tc.wantErr != "" {
				assertErrorContains(t, body, tc.wantErr)
			}
		})
	}

	if got := countRows(t, "addresses"); got != 0 {
		t.Fatalf("addresses count = %d, want 0 (failed requests must not persist anything)", got)
	}
}
