package repository

import (
	"cart/internal/cart"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go"
)

type valkeyCartRepository struct {
	conn valkey.Client
}

func (v *valkeyCartRepository) formatKey(userID uuid.UUID) string {
	return fmt.Sprintf("%s%s", cart.CartKeyPrefix, userID)
}

func (v *valkeyCartRepository) GetCart(ctx context.Context, userID uuid.UUID) (*cart.Cart, error) {
	key := v.formatKey(userID)

	data, err := v.conn.Do(ctx, v.conn.B().Get().Key(key).Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return &cart.Cart{
				UserID: userID,
				Items:  []cart.CartItem{},
				Total:  0,
			}, nil
		}

		return nil, fmt.Errorf("failed to get cart from Valkey: %w", err)
	}

	var c cart.Cart
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart data: %w", err)
	}

	return &c, nil
}

func (v *valkeyCartRepository) SaveCart(ctx context.Context, c *cart.Cart) error {
	key := v.formatKey(c.UserID)
	c.UpdatedAt = time.Now()

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marchal cart: %w", err)
	}

	err = v.conn.Do(ctx, v.conn.B().Set().Key(key).Value(string(data)).Ex(cart.CartTTL).Build()).Error()
	if err != nil {
		return fmt.Errorf("failed to save cart to Valkey: %w", err)
	}

	return nil
}

func (v *valkeyCartRepository) AddItem(ctx context.Context, userID uuid.UUID, item cart.CartItem) error {
	c, err := v.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	for i, existing := range c.Items {
		if existing.ProductID == item.ProductID {
			c.Items[i].Quantity += item.Quantity
			c.Total = c.TotalPrice()
			return v.SaveCart(ctx, c)
		}
	}

	c.Items = append(c.Items, item)
	c.Total = c.TotalPrice()

	return v.SaveCart(ctx, c)
}

func (v *valkeyCartRepository) RemoveItem(ctx context.Context, userID, productID uuid.UUID) error {
	c, err := v.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	items := []cart.CartItem{}
	for _, item := range c.Items {
		if item.ProductID != productID {
			items = append(items, item)
		}
	}

	c.Items = items
	c.Total = c.TotalPrice()
	return v.SaveCart(ctx, c)
}

func (v *valkeyCartRepository) UpdateQuantity(ctx context.Context, userID, productID uuid.UUID, q int) error {
	c, err := v.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items[i].Quantity = q
			c.Total = c.TotalPrice()
			return v.SaveCart(ctx, c)
		}
	}

	return nil
}

func (v *valkeyCartRepository) ClearCart(ctx context.Context, userID uuid.UUID) error {
	key := v.formatKey(userID)

	if err := v.conn.Do(ctx, v.conn.B().Del().Key(key).Build()).Error(); err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}

func (v *valkeyCartRepository) GetCartTTL(ctx context.Context, userID uuid.UUID) (time.Duration, error) {
	key := v.formatKey(userID)

	ttl, err := v.conn.Do(ctx, v.conn.B().Ttl().Key(key).Build()).AsInt64()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return -2 * time.Second, nil
		}

		return 0, fmt.Errorf("failed to get TTL: %w", err)
	}

	return time.Duration(ttl) * time.Second, nil
}
