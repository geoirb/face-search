package file

import (
	service "github.com/geoirb/face-search/internal/face-search"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

// GetPath ...
func (m *Mock) GetPath(file service.File) (path string, err error) {
	args := m.Called(file)
	if a, ok := args.Get(0).(string); ok {
		return a, args.Error(1)
	}
	return "", nil
}

// Delete ...
func (m *Mock) Delete(file service.File) (err error) {
	args := m.Called(file)
	return args.Error(0)
}

// GetHash ...
func (m *Mock) GetHash(file service.File) (hash string, err error) {
	args := m.Called(file)
	if a, ok := args.Get(0).(string); ok {
		return a, args.Error(1)
	}
	return "", nil
}
