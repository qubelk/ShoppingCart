package cart

import (
	"cart/internal/product"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Cart struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func New(ctx context.Context, conn string) (*Cart, error) {
	cfg, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, fmt.Errorf("failed parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to check connection to datebase: %w", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS cart (
		id UUID PRIMARY KEY,
		title VARCHAR(100) NOT NULL,
		description VARCHAR(1000),
		price DECIMAL(10,2) CHECK (price >= 0)
	)
	`

	if _, err := pool.Exec(ctx, createTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create table in database: %w", err)
	}

	return &Cart{ctx: ctx, pool: pool}, err
}

func (c *Cart) Add(p *product.Product) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("failed to validate product: %w", err)
	}

	addQuery := `INSERT INTO cart (id, title, description, price) VALUES ($1, $2, $3, $4)`

	_, err := c.pool.Exec(c.ctx, addQuery, p.ID, p.Title, p.Description, p.Price)
	if err != nil {
		return fmt.Errorf("failed to add product to the cart: %w", err)
	}

	return nil
}

func (c *Cart) GetByID(id uuid.UUID) (*product.Product, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("want non zero ID")
	}

	var p product.Product
	getQuery := `SELECT id, title, description, price FROM cart WHERE id = $1`
	err := c.pool.QueryRow(c.ctx, getQuery, id).Scan(&p.ID, &p.Title, &p.Description, &p.Price)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to get product: product with ID %s not exists", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &p, nil
}

func (c *Cart) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("want non zero ID")
	}

	deleteQuery := `DELETE FROM cart WHERE id = $1`

	tag, err := c.pool.Exec(c.ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product with ID %s not founded", id)
	}

	return nil
}

func (c *Cart) Close() {
	c.pool.Close()
}
