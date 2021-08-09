package chromedp

import (
	"time"

	faceSearch "github.com/geoirb/face-search/internal/face-search"
	"github.com/stretchr/testify/mock"
)

// Mock ...
type Mock struct {
	mock.Mock
}

// FaceSearch ...
func (m *Mock) Face(search faceSearch.Config) ([]byte, error) {
	args := m.Called(search)

	time.Sleep(search.Timeout)

	if res, ok := args.Get(0).([]byte); ok {
		return res, args.Error(1)
	}
	return nil, nil
}
