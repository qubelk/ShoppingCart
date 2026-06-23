package user

import (
	"pkg/testutils"
	"strings"
	"testing"
)

func longString(length int) string {
	var b strings.Builder
	for range length {
		b.WriteString("a")
	}
	return b.String()
}

func TestValidateLogin(t *testing.T) {
	tests := []testutils.TestCase{
		{Name: "valid", Data: "login", WantErr: false},
		{Name: "invalid: too short", Data: "lg", WantErr: true},
		{Name: "invalid: too long", Data: longString(101), WantErr: true},
		{Name: "invalid: empty", Data: "", WantErr: true},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return validateLogin(a.(string))
	})
}

func TestValidateEmail(t *testing.T) {
	tests := []testutils.TestCase{
		{Name: "valid", Data: "example@example.com", WantErr: false},
		{Name: "invalid: wrong format", Data: "example#example,com", WantErr: true},
		{Name: "invalid: too long", Data: longString(256), WantErr: true},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return validateEmail(a.(string))
	})
}

func TestValidatePassword(t *testing.T) {
	tests := []testutils.TestCase{
		{Name: "valid", Data: "ValidPass1", WantErr: false},
		{Name: "invalid: too short", Data: "shortps", WantErr: true},
		{Name: "invalid: too long", Data: longString(256), WantErr: true},
		{Name: "invalid: empty", Data: "", WantErr: true},
		{Name: "invalid: no upper", Data: "lower_case1", WantErr: true},
		{Name: "invalid: no lower", Data: "UPPER_CASE1", WantErr: true},
		{Name: "invalid: no digits", Data: "NoDigits", WantErr: true},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return validatePassword(a.(string))
	})
}

func TestLoginRequestValidate(t *testing.T) {
	tests := []testutils.TestCase{
		{
			Name: "valid login request",
			Data: &LoginRequest{
				Login:    "validuser",
				Password: "ValidPass1",
			},
			WantErr: false,
		},
		{
			Name: "invalid login",
			Data: &LoginRequest{
				Login:    "a",
				Password: "ValidPass1",
			},
			WantErr: true,
		},
		{
			Name: "invalid password",
			Data: &LoginRequest{
				Login:    "validuser",
				Password: "weak",
			},
			WantErr: true,
		},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return a.(*LoginRequest).Validate()
	})
}

func TestRegisterRequestValidate(t *testing.T) {
	tests := []testutils.TestCase{
		{
			Name: "valid register request",
			Data: &RegisterRequest{
				Login:    "newuser",
				Email:    "user@example.com",
				Password: "ValidPass1",
			},
			WantErr: false,
		},
		{
			Name: "invalid login",
			Data: &RegisterRequest{
				Login:    "ab",
				Email:    "user@example.com",
				Password: "ValidPass1",
			},
			WantErr: true,
		},
		{
			Name: "invalid email",
			Data: &RegisterRequest{
				Login:    "newuser",
				Email:    "invalid-email",
				Password: "ValidPass1",
			},
			WantErr: true,
		},
		{
			Name: "invalid password",
			Data: &RegisterRequest{
				Login:    "newuser",
				Email:    "user@example.com",
				Password: "weak",
			},
			WantErr: true,
		},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return a.(*RegisterRequest).Validate()
	})
}
