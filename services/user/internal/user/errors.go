package user

import "errors"

var (
	ErrTooShortPassword = errors.New("password must contains uppercase, lowercase and digit")
	ErrWeakPassword     = errors.New("password must be at least 8 character")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrEmailAlredyExist = errors.New("user with this email already exists")
	ErrUserNotFound     = errors.New("user not found")
)
