package plugin

import (
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) GetExpresionDir() (string, error) {
	args := m.Called()
	if expresionDir, ok := args.Get(0).(string); ok {
		return expresionDir, args.Error(1)
	}
	return "", nil
}
