package cart

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Title     string    `json:"title"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
}

type Cart struct {
	UserID    uuid.UUID  `json:"user_id"`
	Items     []CartItem `json:"items"`
	Total     float64    `json:"total"`
	UpdatedAt time.Time  `json:"updated_at"`
}

const (
	CartKeyPrefix string        = "cart:"
	CartTTL       time.Duration = 7 * 24 * time.Hour
)

func (c *Cart) TotalPrice() float64 {
	total := 0.0
	for _, item := range c.Items {
		total += item.Price * float64(item.Quantity)
	}

	return total
}
