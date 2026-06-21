package service

import (
	"context"
	"fmt"
	"product/internal/product"
	"product/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepository
}

func New(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (ps *ProductService) Create(ctx context.Context, req product.CreateProductRequest) (*product.CreateProductResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("create product request failed validation: %w", err)
	}

	p := product.New(req.Title, req.Description, req.Price)
	if err := ps.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &product.CreateProductResponse{Product: *p}, nil
}

func (ps *ProductService) SearchProducts(ctx context.Context, req *product.SearchProductRequest) (*product.SearchProductsResponse, error) {
	p, err := ps.repo.GetByTitle(ctx, req.Title)
	if err != nil {
		return nil, err
	}

	return &product.SearchProductsResponse{Products: p}, nil
}

func (ps *ProductService) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.GetProductResponse, error) {
	p, err := ps.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &product.GetProductResponse{Product: *p}, nil
}

func (ps *ProductService) Delete(ctx context.Context, req *product.DeleteProductRequest) error {
	return ps.repo.Delete(ctx, req.ID, req.OwnerID)
}
