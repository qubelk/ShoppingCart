package product

import "github.com/google/uuid"

type GetProductRequest struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type SearchProductRequest struct {
	Title string `json:"title"`
}

type SearchProductsResponse struct {
	Products []Product `json:"products"`
}

type GetProductResponse struct {
	Product Product `json:"product"`
}

type CreateProductRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type CreateProductResponse struct {
	Product Product `json:"product"`
}

type DeleteProductRequest struct {
	ID      uuid.UUID `json:"id"`
	OwnerID uuid.UUID `json:"owner_id"`
}
