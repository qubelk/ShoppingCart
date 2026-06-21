package product

import "errors"

var (
	ErrInvalidTitle       = errors.New("product title must contains at least 3 character")
	ErrInvalidDescription = errors.New("description can't be empty")
	ErrInvalidPrice       = errors.New("product price can't be negative")
	ErrProductNotFound    = errors.New("product not found")
	ErrProductNotExists   = errors.New("product not exists")
	ErrInvalidCount       = errors.New("invalid count of product")
)
