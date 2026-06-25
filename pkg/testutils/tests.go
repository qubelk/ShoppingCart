package testutils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name    string
	Data    any
	WantErr bool
	ErrType error
	Setup   func(any)
}

type ServiceTestRunner interface {
	RunTest(ctx context.Context, tc TestCase) (any, error)
}

type RepositoryTestRunner interface {
	RunTest(ctx context.Context, tc TestCase) (any, error)
}

func UnitTestRunner(t *testing.T, tests []TestCase, fn func(any) error) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			if tt.Setup != nil {
				tt.Setup(nil)
			}

			err := fn(tt.Data)
			AssertError(t, err, tt.WantErr, tt.ErrType)
		})
	}
}

func MockTestRunner(t *testing.T, tests []TestCase, runner ServiceTestRunner) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Setup != nil {
				tt.Setup(nil)
			}

			res, err := runner.RunTest(context.Background(), tt)
			AssertError(t, err, tt.WantErr, tt.ErrType)

			if !tt.WantErr && res != nil {
				assert.NotNil(t, res)
			}
		})
	}
}

func AssertError(t *testing.T, err error, wantErr bool, errType error) {
	t.Helper()
	if wantErr {
		require.Error(t, err)
		if errType != nil {
			require.ErrorIs(t, err, errType)
		}
	} else {
		require.NoError(t, err)
	}
}

func GetCurrentTime() time.Time {
	return time.Now().Truncate(time.Millisecond)
}
