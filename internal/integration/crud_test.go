//go:build integration

package integration

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/url"
	"testing"
)

func TestProductImageRoundtrip(t *testing.T) {
	resetDB(t)

	supplierID := createSupplier(t, uniqueName("supplier"), addressBody("Russia", "Moscow", "Image Street"))
	raw := []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00, 0x01}

	productID := createProduct(t, uniqueName("product"), supplierID, 10, map[string]any{
		"jpg_base64": base64.StdEncoding.EncodeToString(raw),
	})

	if got := countRows(t, "images"); got != 1 {
		t.Fatalf("images count = %d, want 1", got)
	}

	status, body := doJSON(t, http.MethodGet, "/products/"+productID+"/image", nil)
	if status != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", status, http.StatusOK, body)
	}
	if !bytes.Equal([]byte(body), raw) {
		t.Fatalf("image bytes = %v, want %v", []byte(body), raw)
	}
}

func TestClientsPaginationAndFilter(t *testing.T) {
	resetDB(t)

	first := uniqueName("client")
	second := uniqueName("client")
	third := uniqueName("client")

	createClient(t, first, "1990-05-10", addressBody("Russia", "Moscow", "Pagination 1"))
	createClient(t, second, "1991-06-11", addressBody("Russia", "Moscow", "Pagination 2"))
	createClient(t, third, "1992-07-12", addressBody("Russia", "Moscow", "Pagination 3"))

	body := mustStatus(t, http.MethodGet, "/clients?limit=2", nil, http.StatusOK)
	var page clientListResponse
	decode(t, body, &page)

	if page.Limit != 2 || page.Offset != 0 {
		t.Fatalf("limit/offset = %d/%d, want 2/0", page.Limit, page.Offset)
	}
	if len(page.Data) != 2 {
		t.Fatalf("data length = %d, want 2", len(page.Data))
	}

	body = mustStatus(t, http.MethodGet, "/clients?limit=2&offset=2", nil, http.StatusOK)
	var nextPage clientListResponse
	decode(t, body, &nextPage)

	if nextPage.Offset != 2 {
		t.Fatalf("offset = %d, want 2", nextPage.Offset)
	}
	if len(nextPage.Data) != 1 {
		t.Fatalf("data length = %d, want 1", len(nextPage.Data))
	}

	body = mustStatus(t, http.MethodGet, "/clients?name="+url.QueryEscape(second), nil, http.StatusOK)
	var filtered clientListResponse
	decode(t, body, &filtered)

	if filtered.Total != 1 || len(filtered.Data) != 1 {
		t.Fatalf("filtered total/data = %d/%d, want 1/1 (body %s)", filtered.Total, len(filtered.Data), body)
	}
	if filtered.Data[0].Name != second {
		t.Fatalf("filtered name = %q, want %q", filtered.Data[0].Name, second)
	}
}

func TestClientDeleteLifecycle(t *testing.T) {
	resetDB(t)

	name := uniqueName("client")
	clientID := createClient(t, name, "1990-05-10", addressBody("Russia", "Moscow", "Delete Street"))

	mustStatus(t, http.MethodDelete, "/clients/"+clientID, nil, http.StatusNoContent)
	if got := countRows(t, "clients"); got != 0 {
		t.Fatalf("clients count = %d, want 0", got)
	}

	status, body := doJSON(t, http.MethodDelete, "/clients/"+clientID, nil)
	if status != http.StatusNotFound {
		t.Fatalf("status = %d, want %d, body = %s", status, http.StatusNotFound, body)
	}
	assertErrorContains(t, body, "client not found")
}
