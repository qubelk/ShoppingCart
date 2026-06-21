package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Login     string    `json:"login" db:"login"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func New(email string, pass string, login string) *User {
	return &User{
		ID:       uuid.New(),
		Email:    email,
		Password: pass,
		Login:    login,
	}
}
