package repository

import (
	"context"
	"product/internal/product"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService interface {
	Create(ctx context.Context, p *product.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error)
	GetByTitle(ctx context.Context, title string) ([]product.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

func New(pool *pgxpool.Pool) ProductService {
	return &pgProductRepository{pool: pool}
}
