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
func (m *Mock) Save(ctx context.Context, fs service.FaceSearch) error {
	return m.Called(fs).Error(0)
}

// Get ...
func (m *Mock) Get(ctx context.Context, filter service.FaceSearchFilter) (service.FaceSearch, error) {
	args := m.Called(filter)

	if res, ok := args.Get(0).(service.FaceSearch); ok {
		return res, args.Error(0)
	}
	return service.FaceSearch{}, nil
}
