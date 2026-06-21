package service

import (
	"context"
	"errors"
	"fmt"
	"user/auth"
	"user/internal/repository"
	"user/internal/user"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo   repository.UserRepository
	jwtKey string
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

	existingLogin, err := s.repo.GetByLogin(ctx, req.Login)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to check login uniqueness: %w", err)
	}

	if existingLogin != nil {
		return nil, user.ErrLoginAlreadyExist
	}

	existingEmail, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
	}

	if existingEmail != nil {
		return nil, user.ErrEmailAlredyExist
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	u := user.New(req.Email, string(hashed), req.Login)
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return &user.RegisterResponse{User: *u}, nil
}

func (s *UserService) GetProfile(ctx context.Context, login string) (*user.UserResponse, error) {
	u, err := s.repo.GetByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user.UserResponse{ID: u.ID, Login: u.Login, CreatedAt: u.CreatedAt}, nil
}

func (s *UserService) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate login request: %w", err)
	}

	u, err := s.repo.GetByLogin(ctx, req.Login)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &user.LoginResponse{User: *u}, nil
}

func (s *UserService) Delete(ctx context.Context, req *user.DeleteRequest) error {
	u, err := s.repo.GetByLogin(ctx, req.Login)
	if err != nil {
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

func (s *UserService) GenerateToken(id string) (string, error) {
	return auth.New().GenerateJWT(id)
}
