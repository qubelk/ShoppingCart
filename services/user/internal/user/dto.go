package user

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type DeleteRequest struct {
	ID       uuid.UUID `json:"id" binding:"required"`
	Login    string    `json:"login" binding:"required"`
	Password string    `json:"password" binding:"required"`
}

type RegisterResponse struct {
	User User `json:"user"`
}

type LoginResponse struct {
	User User `json:"user"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
}
