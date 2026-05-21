// Package producttest provides test helpers for the product domain.
// Import only in _test.go files.
package producttest

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/your-org/your-service/internal/domain/product"
)

// MemoryRepository is an in-memory implementation of product.Repository for tests.
type MemoryRepository struct {
	mu        sync.RWMutex
	products  map[uuid.UUID]*product.Product
	SaveErr   error
	FindErr   error
	ListErr   error
	UpdateErr error
	DeleteErr error
}

// NewMemoryRepository returns an empty MemoryRepository ready for use in tests.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{products: make(map[uuid.UUID]*product.Product)}
}

// Save persists or overwrites the product in the in-memory store.
func (r *MemoryRepository) Save(_ context.Context, p *product.Product) error {
	if r.SaveErr != nil {
		return r.SaveErr
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.products[p.ID()] = p
	return nil
}

// FindByID returns the product with the given ID or ErrNotFound.
func (r *MemoryRepository) FindByID(_ context.Context, id uuid.UUID) (*product.Product, error) {
	if r.FindErr != nil {
		return nil, r.FindErr
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.products[id]
	if !ok {
		return nil, product.ErrNotFound
	}
	return p, nil
}

// List returns a paginated slice of all products in the store.
func (r *MemoryRepository) List(_ context.Context, f product.ListFilter) (product.Page, error) {
	if r.ListErr != nil {
		return product.Page{}, r.ListErr
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := make([]*product.Product, 0, len(r.products))
	for _, p := range r.products {
		all = append(all, p)
	}
	total := int64(len(all))
	start := int(f.Offset)
	if start > len(all) {
		start = len(all)
	}
	end := start + int(f.Limit)
	if end > len(all) {
		end = len(all)
	}
	return product.Page{Items: all[start:end], Total: total}, nil
}

// Update replaces an existing product or returns ErrNotFound.
func (r *MemoryRepository) Update(_ context.Context, p *product.Product) error {
	if r.UpdateErr != nil {
		return r.UpdateErr
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.products[p.ID()]; !ok {
		return product.ErrNotFound
	}
	r.products[p.ID()] = p
	return nil
}

// Delete removes the product with the given ID or returns ErrNotFound.
func (r *MemoryRepository) Delete(_ context.Context, id uuid.UUID) error {
	if r.DeleteErr != nil {
		return r.DeleteErr
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.products[id]; !ok {
		return product.ErrNotFound
	}
	delete(r.products, id)
	return nil
}

// ExistsBySKU reports whether any product in the store has the given SKU.
func (r *MemoryRepository) ExistsBySKU(_ context.Context, sku product.SKU) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.products {
		if p.SKU().String() == sku.String() {
			return true, nil
		}
	}
	return false, nil
}
