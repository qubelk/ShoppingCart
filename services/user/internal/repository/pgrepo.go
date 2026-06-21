package repository

import (
	"context"
	"fmt"
	"user/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgUserRepository struct {
	pool *pgxpool.Pool
}

func (p *pgUserRepository) Create(ctx context.Context, user *user.User) error {
	createQuery := `
		INSERT INTO users (id, email, password, login)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`

	err := p.pool.QueryRow(
		ctx,
		createQuery,
		user.ID,
		user.Email,
		user.Password,
		user.Login,
	).Scan(&user.CreatedAt)

	return err
}

func (p *pgUserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	getQuery := `
		SELECT id, email, password, login, created_at FROM users WHERE email = $1
	`

	var u user.User
	err := p.pool.QueryRow(
		ctx,
		getQuery,
		email,
	).Scan(&u.ID, &u.Email, &u.Password, &u.Login, &u.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &u, nil
}

func (p *pgUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	getQuery := `
		SELECT id, email, password, login, created_at FROM users WHERE id = $1
	`

	var u user.User
	err := p.pool.QueryRow(
		ctx,
		getQuery,
		id,
	).Scan(&u.ID, &u.Email, &u.Password, &u.Login, &u.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &u, nil
}

func (p *pgUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	deleteQuery := `DELETE FROM users WHERE id = $1`

	tag, err := p.pool.Exec(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user with ID %s not found", id)
	}

	return nil
}
