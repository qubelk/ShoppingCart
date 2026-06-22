package user

import (
	"net/mail"
	"regexp"

	"golang.org/x/sync/errgroup"
)

func validateLogin(login string) error {
	if len(login) < 3 || len(login) > 100 {
		return ErrInvalidLogin
	}

	return nil
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)

	if err != nil || len(email) > 255 {
		return ErrInvalidEmail
	}

	return nil
}

func validatePassword(pass string) error {
	if len(pass) < 8 || len(pass) > 255 {
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

	return g.Wait()
}

func (r *DeleteRequest) Validate() error {
	var g errgroup.Group

	g.Go(func() error {
		return validateLogin(r.Login)
	})

	g.Go(func() error {
		return validatePassword(r.Password)
	})

	return g.Wait()
}
