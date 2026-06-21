package user

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Login    string `json:"login"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type DeleteRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	User User `json:"user"`
}

type LoginResponse struct {
	User User `json:"user"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
}
