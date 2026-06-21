package product

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"descrption"`
	Price       float64   `json:"price" db:"price"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Stock       int       `json:"stock" db:"stock"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

func New(title, description string, price float64) *Product {
	return &Product{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Price:       price,
	}
}
