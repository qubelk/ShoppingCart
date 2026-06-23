package cart

import (
	"time"

	"github.com/google/uuid"
)

type AddItemRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Item   CartItem  `json:"item" binding:"required"`
}

type UpdateQuantityRequest struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required"`
}

type RemoveItemRequest struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	ProductID uuid.UUID `json:"product_id" binding:"required"`
}

type ClearCartRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

type GetCartRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

type GetCartTTLRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

type AddItemResponse struct {
	Cart Cart `json:"cart"`
}

type UpdateQuantityResponse struct {
	Cart Cart `json:"cart"`
}

type RemoveItemResponse struct {
	Cart Cart `json:"cart"`
}

type GetCartResponse struct {
	Cart Cart `json:"cart"`
}

type GetCartTTLResponse struct {
	Ttl time.Duration `json:"ttl"`
}
