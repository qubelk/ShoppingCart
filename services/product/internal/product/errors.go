package product

import "errors"

var (
	ErrInvalidTitle       = errors.New("title must contains at least 3 character and no more than 255 characters")
	ErrInvalidDescription = errors.New("description can't be empty or more than 1000 characters")
	ErrInvalidPrice       = errors.New("product price can't be negative")
	ErrProductNotFound    = errors.New("product not found")
	ErrProductNotExists   = errors.New("product not exists")
	ErrInvalidCount       = errors.New("invalid count of product")
)
