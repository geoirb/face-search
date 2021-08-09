package service

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Mock ..
type Mock struct {
	mock.Mock
}

// GetSearchConfig ...
func (m *Mock) GetSearchConfig(ctx context.Context) (cfg Config, err error) {
	args := m.Called()
	if a, ok := args.Get(0).(Config); ok {
		return a, args.Error(1)
	}
	return Config{}, nil
}

// UpdateSearchConfig ...
func (m *Mock) UpdateSearchConfig(ctx context.Context, newSearch Config) error {
	args := m.Called(newSearch)
	return args.Error(0)

}

// FaceSearch ...
func (m *Mock) FaceSearch(ctx context.Context, sfs Search) (Result, error) {
	args := m.Called(sfs)
	if a, ok := args.Get(0).(Result); ok {
		return a, args.Error(1)
	}
	return Result{}, nil
}

// GetFaceSearchResult ...
func (m *Mock) GetFaceSearchResult(ctx context.Context, tfs TaskFaceSearch) (Result, error) {
	args := m.Called(tfs)
	if a, ok := args.Get(0).(Result); ok {
		return a, args.Error(1)
	}
	return Result{}, nil
}
