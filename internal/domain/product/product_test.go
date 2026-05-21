package product_test

import (
	"testing"
	"time"

	"github.com/your-org/your-service/internal/domain/product"
)

func TestNewName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"valid", "My Product", nil},
		{"min length", "A", nil},
		{"max length", string(make([]byte, 100)), nil},
		{"empty", "", product.ErrInvalidName},
		{"too long", string(make([]byte, 101)), product.ErrInvalidName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := product.NewName(tt.input)
			if err != tt.wantErr {
				t.Errorf("NewName(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestNewDescription(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"valid", "A great product", nil},
		{"empty allowed", "", nil},
		{"max length", string(make([]byte, 500)), nil},
		{"too long", string(make([]byte, 501)), product.ErrInvalidDesc},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := product.NewDescription(tt.input)
			if err != tt.wantErr {
				t.Errorf("NewDescription(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestNewPrice(t *testing.T) {
	tests := []struct {
		name    string
		input   float64
		wantErr error
	}{
		{"valid", 9.99, nil},
		{"minimum", 0.01, nil},
		{"zero", 0, product.ErrInvalidPrice},
		{"negative", -1.0, product.ErrInvalidPrice},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := product.NewPrice(tt.input)
			if err != tt.wantErr {
				t.Errorf("NewPrice(%v) error = %v, want %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestNewSKU(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"valid", "SKU-001", nil},
		{"min length", "A", nil},
		{"max length", string(make([]byte, 50)), nil},
		{"empty", "", product.ErrInvalidSKU},
		{"too long", string(make([]byte, 51)), product.ErrInvalidSKU},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := product.NewSKU(tt.input)
			if err != tt.wantErr {
				t.Errorf("NewSKU(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestProduct_MutationMethods(t *testing.T) {
	name, _ := product.NewName("Original")
	desc, _ := product.NewDescription("Original desc")
	price, _ := product.NewPrice(10.0)
	sku, _ := product.NewSKU("SKU-001")
	p := product.New(name, desc, price, sku)

	t.Run("Rename updates name and updatedAt", func(t *testing.T) {
		before := p.UpdatedAt()
		time.Sleep(2 * time.Millisecond) // ensure time advances on low-resolution clocks
		newName, _ := product.NewName("Updated")
		p.Rename(newName)
		if p.Name().String() != "Updated" {
			t.Error("expected name to be Updated")
		}
		if !p.UpdatedAt().After(before) {
			t.Error("expected updatedAt to be refreshed")
		}
	})

	t.Run("Reprice updates price", func(t *testing.T) {
		newPrice, _ := product.NewPrice(99.99)
		p.Reprice(newPrice)
		if p.Price().Value() != 99.99 {
			t.Errorf("expected price 99.99, got %v", p.Price().Value())
		}
	})
}
