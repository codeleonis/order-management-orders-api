package product

// Price is a value object for the product price.
// Invariant: value > 0
type Price struct{ value float64 }

// NewPrice validates and constructs a Price value object.
func NewPrice(v float64) (Price, error) {
	if v <= 0 {
		return Price{}, ErrInvalidPrice
	}
	return Price{value: v}, nil
}

// Value returns the underlying price as a float64.
func (p Price) Value() float64 { return p.value }
