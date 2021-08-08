package mongo

import (
	"context"

	"github.com/stretchr/testify/mock"

	service "github.com/geoirb/face-search/internal/face-search"
)

// Mock ...
type Mock struct {
	mock.Mock
}

// Save ...
func (m *Mock) Save(ctx context.Context, fs service.Result) error {
	return m.Called(fs).Error(0)
}

// Get ...
func (m *Mock) Get(ctx context.Context, filter service.FaceSearchFilter) (service.Result, error) {
	args := m.Called(filter)

	if res, ok := args.Get(0).(service.Result); ok {
		return res, args.Error(1)
	}
	return service.Result{}, nil
}
