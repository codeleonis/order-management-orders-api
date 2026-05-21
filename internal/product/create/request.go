package create

// Request is the JSON body for the create-product endpoint.
type Request struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
}
