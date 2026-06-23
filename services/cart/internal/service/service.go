package service

import (
	"cart/internal/cart"
	"cart/internal/repository"
	"context"
	"fmt"
)

type CartService struct {
	repo repository.CartRepository
}

func New(repo repository.CartRepository) *CartService {
	return &CartService{
		repo: repo,
	}
}

func (cs *CartService) GetCart(ctx context.Context, req *cart.GetCartRequest) (*cart.GetCartResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate get cart request: %w", err)
	}

	c, err := cs.repo.GetCart(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return &cart.GetCartResponse{Cart: *c}, err
}

func (cs *CartService) AddItem(ctx context.Context, req *cart.AddItemRequest) (*cart.AddItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate add item to cart request: %w", err)
	}

	if err := cs.repo.AddItem(ctx, req.UserID, req.Item); err != nil {
		return nil, fmt.Errorf("failed to add item in a cart: %w", err)
	}

	c, err := cs.repo.GetCart(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return &cart.AddItemResponse{Cart: *c}, nil
}

func (cs *CartService) UpdateQuantity(ctx context.Context, req *cart.UpdateQuantityRequest) (*cart.UpdateQuantityResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate update quantity request: %w", err)
	}

	if err := cs.repo.UpdateQuantity(ctx, req.UserID, req.ProductID, req.Quantity); err != nil {
		return nil, fmt.Errorf("failed to update quantity: %w", err)
	}

	c, err := cs.repo.GetCart(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return &cart.UpdateQuantityResponse{Cart: *c}, nil
}

func (cs *CartService) RemoveItem(ctx context.Context, req *cart.RemoveItemRequest) (*cart.RemoveItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate remove item request: %w", err)
	}

	if err := cs.repo.RemoveItem(ctx, req.UserID, req.ProductID); err != nil {
		return nil, fmt.Errorf("failed to remove item from cart: %w", err)
	}

	c, err := cs.repo.GetCart(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return &cart.RemoveItemResponse{Cart: *c}, nil
}

func (cs *CartService) CleanCart(ctx context.Context, req *cart.ClearCartRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("failed to validate clean cart request: %w", err)
	}

	return cs.repo.ClearCart(ctx, req.UserID)
}

func (cs *CartService) GetCartTTL(ctx context.Context, req *cart.GetCartTTLRequest) (*cart.GetCartTTLResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate get cart TTL request: %w", err)
	}

	ttl, err := cs.repo.GetCartTTL(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart TTL: %w", err)
	}

	return &cart.GetCartTTLResponse{Ttl: ttl}, nil
}
