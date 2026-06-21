package repository

import (
	"context"
	"product/internal/product"
)

type ProductService interface {
	Create(ctx context.Context, p *product.Product) error
	GetByID(ctx context.Context, id int) (*product.Product, error)
	GetByTitle(ctx context.Context, title string) ([]product.Product, error)
	Delete(ctx context.Context, id int) error
}
