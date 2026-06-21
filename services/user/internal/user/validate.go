package user

import (
	"net/mail"
	"regexp"

	"golang.org/x/sync/errgroup"
)

type Request interface {
	LoginRequest | RegisterRequest
}

func validateLogin(login string) error {
	if len(login) < 3 {
		return ErrInvalidLogin
	}

	return nil
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
		return ErrTooShortPassword
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
	var g errgroup.Group

	g.Go(func() error {
		return validateLogin(r.Login)
	})

	g.Go(func() error {
		return validatePassword(r.Password)
	})

	return g.Wait()
}

func (r *RegisterRequest) Validate() error {
	var g errgroup.Group

	g.Go(func() error {
		return validateLogin(r.Login)
	})

	g.Go(func() error {
		return validateEmail(r.Email)
	})

	g.Go(func() error {
		return validatePassword(r.Password)
	})

	return nil
}
