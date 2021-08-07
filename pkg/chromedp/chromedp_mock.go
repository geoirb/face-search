package chromedp

import (
	"github.com/stretchr/testify/mock"

	service "github.com/geoirb/face-search/pkg/face-search"
)

// MockChromedp ...
type MockChromedp struct {
	mock.Mock
}

// FaceSearch ...
func (m *MockChromedp) FaceSearch(actions []service.Action) ([]byte, error) {
	args := m.Called(actions)

	if res, ok := args.Get(0).([]byte); ok {
		return res, args.Error(0)
	}
	return nil, nil
}
