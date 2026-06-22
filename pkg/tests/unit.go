package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name    string
	Data    any
	WantErr bool
}

func UnitTest(t *testing.T, tests []TestCase, fn func(any) error) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			err := fn(tt.Data)
			if (err != nil) != tt.WantErr {
				require.Error(t, err, "Unexpected error: %v", err)
			} else if (err == nil) == tt.WantErr {
				require.Error(t, err, "Expected error but not got")
			}
		})
	}
}
