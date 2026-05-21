package update

import "time"

// Response is the JSON body returned after a successful product update.
type Response struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	SKU         string    `json:"sku"`
	UpdatedAt   time.Time `json:"updated_at"`
}
