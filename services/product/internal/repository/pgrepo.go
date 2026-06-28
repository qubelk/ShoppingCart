package repository

import (
	"context"
	"errors"
	"fmt"
	"product/internal/product"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgProductRepository struct {
	pool *pgxpool.Pool
}

func (pg *pgProductRepository) AddToStock(ctx context.Context, id uuid.UUID, count int) error {
	updateQuery := `
		UPDATE products
		SET stock = stock + $1
		WHERE id = $2
	`

	_, err := pg.pool.Exec(
		ctx,
		updateQuery,
		count,
		id,
	)

	return err
}

func (pg *pgProductRepository) GetFromStock(ctx context.Context, id uuid.UUID, count int) error {
	updateQuery := `
		UPDATE products
		SET stock = stock - $1
		WHERE id = $2
			AND stock >= $1;
	`

	tag, err := pg.pool.Exec(ctx, updateQuery, id)
	if err != nil {
		return fmt.Errorf("failed to get products with ID %s: %w", id, err)
	}

	if tag.RowsAffected() == 0 {
		return product.ErrProductNotFound
	}

	return nil
}

func (pg *pgProductRepository) Create(ctx context.Context, p *product.Product) error {
	createQuery := `
		INSERT INTO products (id, owner_id, title, description, price, stock)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at
	`

	err := pg.AddToStock(ctx, p.ID, 1)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	err = pg.pool.QueryRow(
		ctx,
		createQuery,
		p.ID,
		p.OwnerID,
		p.Title,
		p.Description,
		p.Price,
		p.Stock,
	).Scan(&p.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (pg *pgProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	getQuery := `
		SELECT (id, owner_id, title, description, price, stock, created_at) FROM products WHERE id = $1
	`

	var p product.Product
	err := pg.pool.QueryRow(
		ctx,
		getQuery,
		id,
	).Scan(&p.ID, &p.OwnerID, &p.Title, &p.Description, &p.Price, &p.Stock, &p.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product.ErrProductNotExists
		}

		return nil, fmt.Errorf("failed to get product with ID %s: %w", id, err)
	}

	return &p, nil
}

func (pg *pgProductRepository) GetByTitle(ctx context.Context, title string) ([]product.Product, error) {
	getQuery := `
		SELECT (id, owner_id, title, description, price, stock, created_at) FROM products WHERE title = $1
	`

	rows, err := pg.pool.Query(
		ctx,
		getQuery,
		title,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product.ErrProductNotExists
		}
		return nil, fmt.Errorf("failed to get products with title %s: %w", title, err)
	}
	defer rows.Close()

	var products []product.Product
	for rows.Next() {
		var p product.Product
		rows.Scan(&p.ID, &p.OwnerID, &p.Title, &p.Description, &p.Price, &p.Stock, &p.CreatedAt)
		products = append(products, p)
	}

	return products, nil
}

func (pg *pgProductRepository) Delete(ctx context.Context, id, ownerID uuid.UUID) error {
	deleteQuery := `
		DELETE FROM products WHERE id = $1 AND owner_id = $2
	`

	tag, err := pg.pool.Exec(ctx, deleteQuery, id, ownerID)
	if err != nil {
		return fmt.Errorf("failed to delete product with ID %s: %w", id, err)
	}

	if tag.RowsAffected() == 0 {
		return product.ErrProductNotFound
	}

	return nil
}
