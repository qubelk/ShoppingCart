package repository

import (
	"cart/internal/cart"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go"
)

type CartRepository interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*cart.Cart, error)
	SaveCart(ctx context.Context, c *cart.Cart) error
	AddItem(ctx context.Context, userID uuid.UUID, item cart.CartItem) error
	RemoveItem(ctx context.Context, userID, productID uuid.UUID) error
	UpdateQuantity(ctx context.Context, userID, productID uuid.UUID, q int) error
	ClearCart(ctx context.Context, userID uuid.UUID) error
	GetCartTTL(ctx context.Context, userID uuid.UUID) (time.Duration, error)
}

func New(conn valkey.Client) CartRepository {
	return &valkeyCartRepository{
		conn: conn,
	}
}
