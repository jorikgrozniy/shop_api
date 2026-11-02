//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"shop_api/internal/controller/restapi/v1/dto"
)

type rawBody string

type clientListResponse struct {
	Data []struct {
		ID   string `json:"client_id"`
		Name string `json:"name"`
	} `json:"data"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type supplierListResponse struct {
	Data []struct {
		ID   string `json:"supplier_id"`
		Name string `json:"name"`
	} `json:"data"`
}

type productListResponse struct {
	Data []struct {
		ID             string `json:"product_id"`
		Name           string `json:"name"`
		AvailableStock int    `json:"available_stock"`
	} `json:"data"`
}

type productResponse struct {
	ID             string `json:"product_id"`
	Name           string `json:"name"`
	AvailableStock int    `json:"available_stock"`
	ImageID        string `json:"image_id"`
}

var nameSeq atomic.Int64

func uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, nameSeq.Add(1))
}

func resetDB(t *testing.T) {
	t.Helper()
	_, err := pool.Exec(context.Background(), "TRUNCATE clients, suppliers, products, images, addresses CASCADE")
	if err != nil {
		t.Fatalf("reset db: %v", err)
	}
}

func countRows(t *testing.T, table string) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(), "SELECT count(*) FROM "+table).Scan(&n); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return n
}

func sendJSON(method, path string, body any) (int, string, error) {
	var reader io.Reader
	if body != nil {
		switch b := body.(type) {
		case rawBody:
			reader = strings.NewReader(string(b))
		default:
			raw, err := json.Marshal(b)
			if err != nil {
				return 0, "", fmt.Errorf("marshal body: %w", err)
			}
			reader = bytes.NewReader(raw)
		}
	}

	req, err := http.NewRequest(method, baseURL+path, reader)
	if err != nil {
		return 0, "", fmt.Errorf("new request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("do request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", fmt.Errorf("read response: %w", err)
	}
	return resp.StatusCode, string(respBody), nil
}

func doJSON(t *testing.T, method, path string, body any) (int, string) {
	t.Helper()
	status, respBody, err := sendJSON(method, path, body)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	return status, respBody
}

func mustStatus(t *testing.T, method, path string, body any, wantStatus int) string {
	t.Helper()
	status, respBody := doJSON(t, method, path, body)
	if status != wantStatus {
		t.Fatalf("%s %s: status = %d, want %d, body = %s", method, path, status, wantStatus, respBody)
	}
	return respBody
}

func decode(t *testing.T, body string, into any) {
	t.Helper()
	if err := json.Unmarshal([]byte(body), into); err != nil {
		t.Fatalf("decode %q: %v", body, err)
	}
}

func assertErrorContains(t *testing.T, body, want string) {
	t.Helper()
	var resp dto.ErrorResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("decode error response %q: %v", body, err)
	}
	if !strings.Contains(resp.Error, want) {
		t.Fatalf("error = %q, want to contain %q", resp.Error, want)
	}
}

func addressBody(country, city, street string) map[string]any {
	return map[string]any{"country": country, "city": city, "street": street}
}

func clientBody(name, birthdate string, address map[string]any) map[string]any {
	return map[string]any{
		"name":      name,
		"surname":   "Testov",
		"birthdate": birthdate,
		"gender":    "male",
		"address":   address,
	}
}

func createClient(t *testing.T, name, birthdate string, address map[string]any) string {
	t.Helper()
	mustStatus(t, http.MethodPost, "/clients", clientBody(name, birthdate, address), http.StatusCreated)
	return findClientID(t, name)
}

func findClientID(t *testing.T, name string) string {
	t.Helper()
	body := mustStatus(t, http.MethodGet, "/clients?name="+url.QueryEscape(name), nil, http.StatusOK)
	var resp clientListResponse
	decode(t, body, &resp)
	for _, c := range resp.Data {
		if c.Name == name {
			return c.ID
		}
	}
	t.Fatalf("client %q not found in %s", name, body)
	return ""
}

func createSupplier(t *testing.T, name string, address map[string]any) string {
	t.Helper()
	mustStatus(t, http.MethodPost, "/suppliers", map[string]any{
		"name":         name,
		"phone_number": "+79995351675",
		"address":      address,
	}, http.StatusCreated)
	return findSupplierID(t, name)
}

func findSupplierID(t *testing.T, name string) string {
	t.Helper()
	body := mustStatus(t, http.MethodGet, "/suppliers?limit=100", nil, http.StatusOK)
	var resp supplierListResponse
	decode(t, body, &resp)
	for _, s := range resp.Data {
		if s.Name == name {
			return s.ID
		}
	}
	t.Fatalf("supplier %q not found in %s", name, body)
	return ""
}

func createProduct(t *testing.T, name, supplierID string, stock int, image map[string]any) string {
	t.Helper()
	payload := map[string]any{
		"name":            name,
		"category":        "Test",
		"price":           99.9,
		"available_stock": stock,
		"supplier_id":     supplierID,
	}
	if image != nil {
		payload["image"] = image
	}
	mustStatus(t, http.MethodPost, "/products", payload, http.StatusCreated)
	return findProductID(t, name)
}

func findProductID(t *testing.T, name string) string {
	t.Helper()
	body := mustStatus(t, http.MethodGet, "/products?limit=100", nil, http.StatusOK)
	var resp productListResponse
	decode(t, body, &resp)
	for _, p := range resp.Data {
		if p.Name == name {
			return p.ID
		}
	}
	t.Fatalf("product %q not found in %s", name, body)
	return ""
}

func getProduct(t *testing.T, id string) productResponse {
	t.Helper()
	body := mustStatus(t, http.MethodGet, "/products/"+id, nil, http.StatusOK)
	var resp productResponse
	decode(t, body, &resp)
	return resp
}
