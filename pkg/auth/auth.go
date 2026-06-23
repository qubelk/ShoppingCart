package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthorizationService struct {
	jwtKey string
}

func New() *AuthorizationService {
	return &AuthorizationService{
		jwtKey: os.Getenv("JWT_SECRET_KEY"),
	}
}

func (s *AuthorizationService) GenerateJWT(id uuid.UUID, login string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   id.String(),
		"login": login,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 12).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtKey))
}

func (s *AuthorizationService) ValidateJWT(tokenString string) (uuid.UUID, string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(s.jwtKey), nil
	})

	if err != nil {
		return uuid.Nil, "", fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"].(string)
		if !ok {
			return uuid.Nil, "", fmt.Errorf("invalid token claims: missing sub")
		}

		id, err := uuid.Parse(sub)
		if err != nil {
			return uuid.Nil, "", fmt.Errorf("invalid user ID in token: %w", err)
		}

		login, _ := claims["login"].(string)

		return id, login, nil
	}

	return uuid.Nil, "", fmt.Errorf("invalid token")
}
