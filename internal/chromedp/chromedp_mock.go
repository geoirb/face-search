package chromedp

import (
	"time"

	"github.com/stretchr/testify/mock"

	service "github.com/geoirb/face-search/internal/face-search"
)

// Mock ...
type Mock struct {
	mock.Mock
}

// FaceSearch ...
func (m *Mock) Face(search service.SearchConfig) ([]byte, error) {
	args := m.Called(search)

	time.Sleep(search.Timeout)

	if res, ok := args.Get(0).([]byte); ok {
		return res, args.Error(1)
	}
	return nil, nil
}
