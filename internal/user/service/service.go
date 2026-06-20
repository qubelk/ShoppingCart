package service

import (
	"cart/internal/user"
	"cart/internal/user/repository"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func New(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate register request: %w", err)
	}

	existingUser, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
	}

	if existingUser != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	u := user.New(req.Email, string(hashed), req.Name)
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return &user.RegisterResponse{User: *u}, nil
}

func (s *UserService) GetProfile(ctx context.Context, id uuid.UUID) (*user.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not founded")
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return u, nil
}

func (s *UserService) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate login request: %w", err)
	}

	u, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not founded")
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &user.LoginResponse{User: *u}, nil
}

func (s *UserService) Delete(ctx context.Context, req *user.DeleteRequest) error {
	u, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user not founded")
		}

		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return fmt.Errorf("invalid password")
	}

	if err := s.repo.Delete(ctx, u.ID); err != nil {
		return err
	}

	return nil
}
