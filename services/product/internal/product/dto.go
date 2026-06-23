package product

import "github.com/google/uuid"

type GetProductRequest struct {
	ID uuid.UUID `json:"id"`
}

type SearchProductRequest struct {
	Title string `json:"title" from:"title"`
}

type SearchProductsResponse struct {
	Products []Product `json:"products"`
}

type GetProductResponse struct {
	Product Product `json:"product"`
}

type CreateProductRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Count       int     `json:"count" binding:"required"`
}

type CreateProductResponse struct {
	Product Product `json:"product"`
}

type DeleteProductRequest struct {
	ID      uuid.UUID `json:"id" binding:"required"`
	OwnerID uuid.UUID `json:"-"`
}
