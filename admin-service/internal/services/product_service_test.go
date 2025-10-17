package services

import (
	"testing"

	"github.com/gamegear/admin-service/internal/models"
	"github.com/google/go-cmp/cmp"
)

func TestDecodeProduct(t *testing.T) {
	payload := []byte(`{
		"id": 42,
		"name": "Gaming Keyboard",
		"description": "Switch Blue",
		"price": 4500,
		"stock": 30,
		"category_id": 1,
		"image_url": "https://example.com/keyboard.jpg"
	}`)

	tests := []struct {
		name    string
		body    []byte
		want    models.ProductSummary
		wantErr bool
	}{
		{
			name: "flat payload",
			body: payload,
			want: expectedProductSummary(),
		},
		{
			name: "data envelope",
			body: []byte(`{"data": ` + string(payload) + `}`),
			want: expectedProductSummary(),
		},
		{
			name: "nested envelope",
			body: []byte(`{"status":"success","data":{"result":{"data": ` + string(payload) + `}}}`),
			want: expectedProductSummary(),
		},
		{
			name:    "no product data",
			body:    []byte(`{"message":"missing"}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeProduct(tt.body)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("decodeProduct() error = %v", err)
			}

			if diff := cmp.Diff(tt.want, *got); diff != "" {
				t.Fatalf("unexpected product (-want +got):\n%s", diff)
			}
		})
	}
}

func expectedProductSummary() models.ProductSummary {
	return models.ProductSummary{
		ID:          42,
		Name:        "Gaming Keyboard",
		Description: "Switch Blue",
		Price:       4500,
		Stock:       30,
		CategoryID:  1,
		ImageURL:    "https://example.com/keyboard.jpg",
	}
}
