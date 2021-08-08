package parser

import (
	"github.com/stretchr/testify/mock"

	service "github.com/geoirb/face-search/internal/face-search"
)

// Mock ...
type Mock struct {
	mock.Mock
}

// GetProfileList ...
func (m *Mock) GetProfileList(payload []byte) []service.Profile {
	if p, ok := m.Called(payload).Get(0).([]service.Profile); ok {
		return p
	}
	return nil
}
