package product

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Title       string
	Description sql.NullString
	Price       float64
}

func New(title string, description sql.NullString, price float64) *Product {
	return &Product{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Price:       price,
	}
}

func (p *Product) Validate() error {
	if p.Title == "" {
		return fmt.Errorf("product title cannot be empty")
	}

	if p.Price < 0 {
		return fmt.Errorf("product price cannot be lower than 0")
	}

	return nil
}

func (p *Product) String() string {
	return fmt.Sprintf(
		"Product{ID: %s, Title: %s, Description: %s, Price: %f}",
		p.ID,
		p.Title,
		p.Description.String,
		p.Price)
}
