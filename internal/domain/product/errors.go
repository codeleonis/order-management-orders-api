package product

import "errors"

// Sentinel errors for the product domain.
var (
	ErrNotFound     = errors.New("product not found")
	ErrDuplicateSKU = errors.New("product with this SKU already exists")
	ErrInvalidName  = errors.New("name must be between 1 and 100 characters")
	ErrInvalidDesc  = errors.New("description must not exceed 500 characters")
	ErrInvalidPrice = errors.New("price must be greater than zero")
	ErrInvalidSKU   = errors.New("SKU must be between 1 and 50 characters")
)
