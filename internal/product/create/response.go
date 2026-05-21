package create

import "time"

// Response is the JSON body returned after a successful product creation.
type Response struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	SKU         string    `json:"sku"`
	CreatedAt   time.Time `json:"created_at"`
}
