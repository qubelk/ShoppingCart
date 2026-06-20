package service

import (
	"cart/internal/user"
	"cart/internal/user/repository"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
		// JWT_SECRET_KEY environment variable located in .env file.
		jwtKey: os.Getenv("JWT_SECRET_KEY"),
	}
}

func (s *UserService) generateJWT(id string) (string, error) {
	claims := jwt.MapClaims{
		"sub": id,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtKey))
}

func (s *UserService) ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(s.jwtKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["sub"].(string)
		if !ok {
			return "", fmt.Errorf("invalid token claims")
		}

		return id, nil
	}

	return "", fmt.Errorf("invalid token")
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
		return nil, user.ErrEmailAlredyExist
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
			return nil, user.ErrUserNotFound
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
			return nil, user.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &user.LoginResponse{User: *u}, nil
}

func (s *UserService) GenerateToken(id string) (string, error) {
	return s.generateJWT(id)
}

func (s *UserService) Delete(ctx context.Context, req *user.DeleteRequest) error {
	u, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.ErrUserNotFound
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
