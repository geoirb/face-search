package mongo

import (
	"context"

	"github.com/stretchr/testify/mock"

	search "github.com/geoirb/face-search/internal/face-search"
)

// Mock ...
type Mock struct {
	mock.Mock
}

// Save ...
func (m *Mock) Save(ctx context.Context, fs search.Result) error {
	return m.Called(fs).Error(0)
}

// Update ...
func (m *Mock) Update(ctx context.Context, fs search.Result) error {
	return m.Called(fs).Error(0)
}

// Get ...
func (m *Mock) Get(ctx context.Context, filter search.ResultFilter) (search.Result, error) {
	args := m.Called(filter)

	if res, ok := args.Get(0).(search.Result); ok {
		return res, args.Error(1)
	}
	return search.Result{}, nil
}
