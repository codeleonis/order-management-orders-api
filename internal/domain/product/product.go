package product

import (
	"time"

	"github.com/google/uuid"
)

// Product is the aggregate root.
// All mutations go through its methods; fields are unexported to enforce invariants.
type Product struct {
	id          uuid.UUID
	name        Name
	description Description
	price       Price
	sku         SKU
	createdAt   time.Time
	updatedAt   time.Time
}

// New creates a new Product. Validation is performed by each Value Object constructor.
func New(name Name, desc Description, price Price, sku SKU) *Product {
	now := time.Now().UTC()
	return &Product{
		id:          uuid.New(),
		name:        name,
		description: desc,
		price:       price,
		sku:         sku,
		createdAt:   now,
		updatedAt:   now,
	}
}

// Reconstitute rebuilds a Product from persisted data, bypassing invariant checks.
// Use only in repository implementations.
func Reconstitute(id uuid.UUID, name Name, desc Description, price Price, sku SKU, createdAt, updatedAt time.Time) *Product {
	return &Product{
		id:          id,
		name:        name,
		description: desc,
		price:       price,
		sku:         sku,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// ID returns the product's unique identifier.
func (p *Product) ID() uuid.UUID { return p.id }

// Name returns the product's name value object.
func (p *Product) Name() Name { return p.name }

// Description returns the product's description value object.
func (p *Product) Description() Description { return p.description }

// Price returns the product's price value object.
func (p *Product) Price() Price { return p.price }

// SKU returns the product's stock-keeping unit value object.
func (p *Product) SKU() SKU { return p.sku }

// CreatedAt returns the UTC timestamp when the product was created.
func (p *Product) CreatedAt() time.Time { return p.createdAt }

// UpdatedAt returns the UTC timestamp of the last mutation.
func (p *Product) UpdatedAt() time.Time { return p.updatedAt }

// Rename updates the product's name and refreshes updatedAt.
func (p *Product) Rename(name Name) {
	p.name = name
	p.updatedAt = time.Now().UTC()
}

// Reprice updates the product's price and refreshes updatedAt.
func (p *Product) Reprice(price Price) {
	p.price = price
	p.updatedAt = time.Now().UTC()
}

// Redescribe updates the product's description and refreshes updatedAt.
func (p *Product) Redescribe(desc Description) {
	p.description = desc
	p.updatedAt = time.Now().UTC()
}

// UpdateSKU updates the product's SKU and refreshes updatedAt.
func (p *Product) UpdateSKU(sku SKU) {
	p.sku = sku
	p.updatedAt = time.Now().UTC()
}
