package product

// SKU is a value object for the stock-keeping unit identifier.
// Invariant: 1 ≤ len(value) ≤ 50
type SKU struct{ value string }

// NewSKU validates and constructs a SKU value object.
func NewSKU(s string) (SKU, error) {
	if len(s) == 0 || len(s) > 50 {
		return SKU{}, ErrInvalidSKU
	}
	return SKU{value: s}, nil
}

// String returns the underlying SKU string.
func (s SKU) String() string { return s.value }
