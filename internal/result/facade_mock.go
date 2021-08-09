package result

import (
	"context"

	"github.com/stretchr/testify/mock"

	search "github.com/geoirb/face-search/internal/face-search"
)

type Mock struct {
	mock.Mock
}

// Get ...
func (m *Mock) Get(ctx context.Context, filter search.ResultFilter) (search.IResult, error) {
	args := m.Called(filter)
	if a, ok := args.Get(0).(search.IResult); ok {
		return a, args.Error(1)
	}
	return nil, nil
}

// New ...
func (m *Mock) New(ctx context.Context, hash string) error {
	args := m.Called(hash)
	return args.Error(0)
}

// GetStatus ...
func (m *Mock) GetStatus() string {
	args := m.Called()
	if a, ok := args.Get(0).(string); ok {
		return a
	}
	return ""
}

// GetUUID ...
func (m *Mock) GetUUID() string {
	args := m.Called()
	if a, ok := args.Get(0).(string); ok {
		return a
	}
	return ""
}

// GetData ...
func (m *Mock) GetData() search.Result {
	if a, ok := m.Called().Get(0).(search.Result); ok {
		return a
	}
	return search.Result{}
}

// SetInProgress ...
func (m *Mock) SetInProgress(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}

// SetSuccess ...
func (m *Mock) SetSuccess(ctx context.Context, profiles []search.Profile) error {
	args := m.Called(profiles)
	return args.Error(0)
}

// SetFailed ...
func (m *Mock) SetFailed(ctx context.Context, err error) error {
	args := m.Called(err)
	return args.Error(0)
}
