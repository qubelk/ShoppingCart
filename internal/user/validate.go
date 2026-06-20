package user

import (
	"fmt"
	"net/mail"
	"regexp"
)

type Request interface {
	LoginRequest | RegisterRequest
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)

	if err != nil {
		return ErrInvalidEmail
	}

	return nil
}

func validatePassword(pass string) error {
	if len(pass) < 8 {
		return ErrToShortPassword
	}

	hasUpper := regexp.MustCompile("[A-Z]").MatchString(pass)
	hasLower := regexp.MustCompile("[a-z]").MatchString(pass)
	hasDigit := regexp.MustCompile("[0-9]").MatchString(pass)

	if !hasUpper || !hasLower || !hasDigit {
		return ErrWeakPassword
	}

	return nil
}

func (r *LoginRequest) Validate() error {
	if err := validateEmail(r.Email); err != nil {
		return fmt.Errorf("email validation failed: %w", err)
	}

	if err := validatePassword(r.Password); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	return nil
}

func (r *RegisterRequest) Validate() error {
	if err := validateEmail(r.Email); err != nil {
		return fmt.Errorf("email validation failed: %w", err)
	}

	if err := validatePassword(r.Password); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	return nil
}
