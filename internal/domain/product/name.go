package product

// Name is a value object for the product name.
// Invariant: 1 ≤ len(value) ≤ 100
type Name struct{ value string }

// NewName validates and constructs a Name value object.
func NewName(s string) (Name, error) {
	if len(s) == 0 || len(s) > 100 {
		return Name{}, ErrInvalidName
	}
	return Name{value: s}, nil
}

// String returns the underlying name string.
func (n Name) String() string { return n.value }
