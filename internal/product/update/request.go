package update

// Request uses pointer fields so nil means "field not provided — do not update".
type Request struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	SKU         *string  `json:"sku"`
}
