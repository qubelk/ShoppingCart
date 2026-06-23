package testutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockRepository interface {
	AssertExpectations(t *testing.T)
}

type ServiceTestHelper struct {
	t     *testing.T
	ctx   context.Context
	mocks []MockRepository
}

func New(t *testing.T) *ServiceTestHelper {
	return &ServiceTestHelper{
		t:     t,
		ctx:   context.Background(),
		mocks: make([]MockRepository, 0),
	}
}

func (s *ServiceTestHelper) With(m MockRepository) *ServiceTestHelper {
	s.mocks = append(s.mocks, m)
	return s
}

func (s *ServiceTestHelper) RunTest(tc TestCase, fn func() error) {
	s.t.Helper()
	s.t.Run(tc.Name, func(t *testing.T) {
		if tc.Setup != nil {
			tc.Setup(nil)
		}

		err := fn()
		AssertError(t, err, tc.WantErr, tc.ErrType)
	})
}

func (h *ServiceTestHelper) RunServiceTest(tc TestCase, runner ServiceTestRunner) {
	h.t.Helper()
	h.t.Run(tc.Name, func(t *testing.T) {
		if tc.Setup != nil {
			tc.Setup(nil)
		}

		result, err := runner.RunTest(h.ctx, tc)
		AssertError(t, err, tc.WantErr, tc.ErrType)

		if !tc.WantErr && result != nil {
			assert.NotNil(t, result)
		}

		h.AssertMocks()
	})
}

func (h *ServiceTestHelper) AssertMocks() {
	for _, mock := range h.mocks {
		mock.AssertExpectations(h.t)
	}
}
