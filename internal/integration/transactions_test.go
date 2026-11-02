//go:build integration

package integration

import (
	"encoding/base64"
	"net/http"
	"sync"
	"testing"

	"github.com/google/uuid"
)

func TestProductImageRollback(t *testing.T) {
	resetDB(t)

	status, body := doJSON(t, http.MethodPost, "/products/"+uuid.NewString()+"/image", map[string]any{
		"jpg_base64": base64.StdEncoding.EncodeToString([]byte("image-bytes")),
	})

	if status != http.StatusNotFound {
		t.Fatalf("status = %d, want %d, body = %s", status, http.StatusNotFound, body)
	}
	assertErrorContains(t, body, "product not found")

	if got := countRows(t, "images"); got != 0 {
		t.Fatalf("images count = %d, want 0 (image insert must be rolled back)", got)
	}
	if got := countRows(t, "products"); got != 0 {
		t.Fatalf("products count = %d, want 0", got)
	}
}

func TestClientAddressRollback(t *testing.T) {
	resetDB(t)

	status, body := doJSON(t, http.MethodPatch, "/clients/"+uuid.NewString()+"/address",
		addressBody("Russia", "Moscow", "Rollback Street"))

	if status != http.StatusNotFound {
		t.Fatalf("status = %d, want %d, body = %s", status, http.StatusNotFound, body)
	}
	assertErrorContains(t, body, "client not found")

	if got := countRows(t, "addresses"); got != 0 {
		t.Fatalf("addresses count = %d, want 0 (address insert must be rolled back)", got)
	}
}

func TestAddressDedupSequential(t *testing.T) {
	resetDB(t)
	address := addressBody("Russia", "Moscow", "Shared Street")

	first := createSupplier(t, uniqueName("supplier"), address)
	second := createSupplier(t, uniqueName("supplier"), address)

	if first == second {
		t.Fatalf("suppliers share id %q", first)
	}
	if got := countRows(t, "addresses"); got != 1 {
		t.Fatalf("addresses count = %d, want 1 (address must be deduplicated)", got)
	}
	if got := countRows(t, "suppliers"); got != 2 {
		t.Fatalf("suppliers count = %d, want 2", got)
	}
}

func TestAddressDedupConcurrent(t *testing.T) {
	resetDB(t)
	const goroutines = 8
	address := addressBody("Russia", "Moscow", "Concurrent Street")

	names := make([]string, goroutines)
	for i := range names {
		names[i] = uniqueName("client")
	}

	statuses := make([]int, goroutines)
	errs := make([]error, goroutines)

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			statuses[i], _, errs[i] = sendJSON(http.MethodPost, "/clients", clientBody(names[i], "1990-05-10", address))
		}(i)
	}
	wg.Wait()

	for i := 0; i < goroutines; i++ {
		if errs[i] != nil {
			t.Fatalf("request %d: %v", i, errs[i])
		}
		if statuses[i] != http.StatusCreated {
			t.Fatalf("request %d: status = %d, want %d", i, statuses[i], http.StatusCreated)
		}
	}
	if got := countRows(t, "addresses"); got != 1 {
		t.Fatalf("addresses count = %d, want 1 (concurrent dedup)", got)
	}
	if got := countRows(t, "clients"); got != goroutines {
		t.Fatalf("clients count = %d, want %d", got, goroutines)
	}
}
