package repository

import (
	"cart/internal/cart"
	"fmt"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go"
)

type valkeyCartRepository struct {
	conn *valkey.Client
}

func (v *valkeyCartRepository) formatKey(userID uuid.UUID) string {
	return fmt.Sprintf("%s%s", cart.CartKeyPrefix, userID)
}
