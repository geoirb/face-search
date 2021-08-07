package parser

import (
	"github.com/stretchr/testify/mock"

	service "github.com/geoirb/face-search/pkg/face-search"
)

// MockParser ...
type MockParser struct {
	mock.Mock
}

func (m *MockParser) GetProfileList(payload []byte) []service.Profile {
	if p, ok := m.Called(payload).Get(0).([]service.Profile); ok {
		return p
	}
	return nil
}
