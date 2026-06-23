package cart

import "errors"

var (
	ErrCartNotFound    = errors.New("cart not found")
	ErrItemNotFound    = errors.New("item not found in cart")
	ErrInvalidID       = errors.New("product ID is required")
	ErrInvalidTitle    = errors.New("product title must be at least 3 character and not more than 255 characters")
	ErrInvalidPrice    = errors.New("price must be greater than zero")
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")
)
