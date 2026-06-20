package user

import "errors"

var (
	ErrToShortPassword  = errors.New("password must contains uppercase, lowercase and digit")
	ErrWeakPassword     = errors.New("password must be at least 8 character")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrEmailAlredyExist = errors.New("user with this email already exists")
	ErrUserNotFounded   = errors.New("user not founded")
)
