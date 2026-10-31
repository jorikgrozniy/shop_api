//go:build integration

package integration

import (
	"net/http"
	"testing"
)

func TestDecreaseAvailableStock(t *testing.T) {
	const initialStock = 30

	cases := []struct {
		name       string
		amount     int
		wantStatus int
		wantStock  int
		wantErr    string
	}{
		{
			name:       "insufficient stock",
			amount:     initialStock + 1,
			wantStatus: http.StatusBadRequest,
			wantStock:  initialStock,
			wantErr:    "insufficient stock",
		},
		{
			name:       "valid decrease",
			amount:     5,
			wantStatus: http.StatusOK,
			wantStock:  initialStock - 5,
		},
		{
			name:       "zero amount",
			amount:     0,
			wantStatus: http.StatusBadRequest,
			wantStock:  initialStock,
			wantErr:    "invalid JSON format",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resetDB(t)

			supplierID := createSupplier(t, uniqueName("supplier"), addressBody("Russia", "Moscow", "Stock Street"))
			productID := createProduct(t, uniqueName("product"), supplierID, initialStock, nil)

			status, body := doJSON(t, http.MethodPost, "/products/"+productID+"/decrease-stock",
				map[string]any{"amount": tc.amount})

			if status != tc.wantStatus {
				t.Fatalf("status = %d, want %d, body = %s", status, tc.wantStatus, body)
			}
			if tc.wantErr != "" {
				assertErrorContains(t, body, tc.wantErr)
			}
			if got := getProduct(t, productID).AvailableStock; got != tc.wantStock {
				t.Fatalf("available_stock = %d, want %d", got, tc.wantStock)
			}
		})
	}
}
