package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"user/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestServiceRegister(t *testing.T) {
	tests := []struct {
		name    string
		req     *user.RegisterRequest
		setup   func(*MockUserRepository)
		wantErr bool
		errType error
	}{
		{
			name: "successful registration",
			req: &user.RegisterRequest{
				Email:    "test@example.com",
				Password: "ValidPass1",
				Login:    "testuser",
			},
			setup: func(m *MockUserRepository) {
				m.On("GetByLogin", mock.Anything, "testuser").
					Return(nil, pgx.ErrNoRows)
				m.On("GetByEmail", mock.Anything, "test@example.com").
					Return(nil, pgx.ErrNoRows)
				m.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "login already exists",
			req: &user.RegisterRequest{
				Email:    "test@example.com",
				Password: "ValidPass1",
				Login:    "existinguser",
			},
			setup: func(m *MockUserRepository) {
				existingUser := &user.User{
					ID:    uuid.New(),
					Login: "existinguser",
					Email: "old@example.com",
				}
				m.On("GetByLogin", mock.Anything, "existinguser").
					Return(existingUser, nil)
			},
			wantErr: true,
			errType: user.ErrLoginAlreadyExist,
		},
		{
			name: "email already exists",
			req: &user.RegisterRequest{
				Email:    "existing@example.com",
				Password: "ValidPass1",
				Login:    "newuser",
			},
			setup: func(m *MockUserRepository) {
				m.On("GetByLogin", mock.Anything, "newuser").
					Return(nil, pgx.ErrNoRows)
				existingUser := &user.User{
					ID:    uuid.New(),
					Login: "olduser",
					Email: "existing@example.com",
				}
				m.On("GetByEmail", mock.Anything, "existing@example.com").
					Return(existingUser, nil)
			},
			wantErr: true,
			errType: user.ErrEmailAlredyExist,
		},
		{
			name: "database error on GetByLogin",
			req: &user.RegisterRequest{
				Email:    "test@example.com",
				Password: "ValidPass1",
				Login:    "testuser",
			},
			setup: func(m *MockUserRepository) {
				m.On("GetByLogin", mock.Anything, "testuser").
					Return(nil, errors.New("database connection error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setup(mockRepo)

			service := New(mockRepo)
			resp, err := service.Register(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.req.Email, resp.User.Email)
				assert.Equal(t, tt.req.Login, resp.User.Login)
				err := bcrypt.CompareHashAndPassword(
					[]byte(resp.User.Password),
					[]byte(tt.req.Password),
				)
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestServiceLogin(t *testing.T) {
	password := "ValidPass1"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name    string
		req     *user.LoginRequest
		setup   func(*MockUserRepository)
		wantErr bool
		errType error
	}{
		{
			name: "successful login",
			req: &user.LoginRequest{
				Login:    "testuser",
				Password: password,
			},
			setup: func(m *MockUserRepository) {
				u := &user.User{
					ID:       uuid.New(),
					Login:    "testuser",
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}
				m.On("GetByLogin", mock.Anything, "testuser").
					Return(u, nil)
			},
			wantErr: false,
		},
		{
			name: "user not found",
			req: &user.LoginRequest{
				Login:    "nonexistent",
				Password: password,
			},
			setup: func(m *MockUserRepository) {
				m.On("GetByLogin", mock.Anything, "nonexistent").
					Return(nil, user.ErrUserNotFound)
			},
			wantErr: true,
			errType: user.ErrUserNotFound,
		},
		{
			name: "invalid password",
			req: &user.LoginRequest{
				Login:    "testuser",
				Password: "WrongPass1",
			},
			setup: func(m *MockUserRepository) {
				u := &user.User{
					ID:       uuid.New(),
					Login:    "testuser",
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}
				m.On("GetByLogin", mock.Anything, "testuser").
					Return(u, nil)
			},
			wantErr: true,
		},
		{
			name: "invalid request data - short password",
			req: &user.LoginRequest{
				Login:    "testuser",
				Password: "short",
			},
			setup:   func(m *MockUserRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setup(mockRepo)

			service := New(mockRepo)
			resp, err := service.Login(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.req.Login, resp.User.Login)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestServiceGetProfile(t *testing.T) {
	tests := []struct {
		name    string
		login   string
		setup   func(*MockUserRepository)
		wantErr bool
		errType error
	}{
		{
			name:  "successful get profile",
			login: "testuser",
			setup: func(m *MockUserRepository) {
				u := &user.User{
					ID:        uuid.New(),
					Login:     "testuser",
					Email:     "test@example.com",
					CreatedAt: time.Now(),
				}
				m.On("GetByLogin", mock.Anything, "testuser").
					Return(u, nil)
			},
			wantErr: false,
		},
		{
			name:  "user not found",
			login: "nonexistent",
			setup: func(m *MockUserRepository) {
				m.On("GetByLogin", mock.Anything, "nonexistent").
					Return(nil, user.ErrUserNotFound)
			},
			wantErr: true,
			errType: user.ErrUserNotFound,
		},
		{
			name:  "database error",
			login: "testuser",
			setup: func(m *MockUserRepository) {
				m.On("GetByLogin", mock.Anything, "testuser").
					Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setup(mockRepo)

			service := New(mockRepo)
			resp, err := service.GetProfile(context.Background(), tt.login)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.login, resp.Login)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestServiceDelete(t *testing.T) {
	password := "ValidPass1"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name    string
		req     *user.DeleteRequest
		setup   func(*MockUserRepository)
		wantErr bool
		errType error
	}{
		{
			name: "successful delete",
			req: &user.DeleteRequest{
				Login:    "testuser",
				Password: password,
			},
			setup: func(m *MockUserRepository) {
				u := &user.User{
					ID:       uuid.New(),
					Login:    "testuser",
					Password: string(hashedPassword),
				}
				m.On("GetByLogin", mock.Anything, "testuser").
					Return(u, nil)
				m.On("Delete", mock.Anything, u.ID).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "user not found",
			req: &user.DeleteRequest{
				Login:    "nonexistent",
				Password: password,
			},
			setup: func(m *MockUserRepository) {
				m.On("GetByLogin", mock.Anything, "nonexistent").
					Return(nil, user.ErrUserNotFound)
			},
			wantErr: true,
			errType: user.ErrUserNotFound,
		},
		{
			name: "invalid password",
			req: &user.DeleteRequest{
				Login:    "testuser",
				Password: "WrongPass1",
			},
			setup: func(m *MockUserRepository) {
				u := &user.User{
					ID:       uuid.New(),
					Login:    "testuser",
					Password: string(hashedPassword),
				}
				m.On("GetByLogin", mock.Anything, "testuser").
					Return(u, nil)
			},
			wantErr: true,
		},
		{
			name: "delete error",
			req: &user.DeleteRequest{
				Login:    "testuser",
				Password: password,
			},
			setup: func(m *MockUserRepository) {
				u := &user.User{
					ID:       uuid.New(),
					Login:    "testuser",
					Password: string(hashedPassword),
				}
				m.On("GetByLogin", mock.Anything, "testuser").
					Return(u, nil)
				m.On("Delete", mock.Anything, u.ID).
					Return(errors.New("delete failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setup(mockRepo)

			service := New(mockRepo)
			err := service.Delete(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestServiceGenerateToken(t *testing.T) {
	service := New(&MockUserRepository{})

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "valid id",
			id:      uuid.New().String(),
			wantErr: false,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.GenerateToken(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}
