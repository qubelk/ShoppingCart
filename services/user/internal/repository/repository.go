package repository

import (
	"context"
	"user/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *user.User) error
	GetByLogin(ctx context.Context, login string) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

func New(pool *pgxpool.Pool) UserRepository {
	return &pgUserRepository{pool: pool}
}
