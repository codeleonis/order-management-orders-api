package product

// Description is a value object for the product description.
// Invariant: len(value) ≤ 500 (empty is allowed)
type Description struct{ value string }

// NewDescription validates and constructs a Description value object.
func NewDescription(s string) (Description, error) {
	if len(s) > 500 {
		return Description{}, ErrInvalidDesc
	}
	return Description{value: s}, nil
}

// String returns the underlying description string.
func (d Description) String() string { return d.value }
