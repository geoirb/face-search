package file

import (
	"testing"

	search "github.com/geoirb/face-search/internal/face-search"
	"github.com/stretchr/testify/assert"
)

func TestGetPath(t *testing.T) {
	f := NewFacade("")
	_, err := f.GetPath(search.File{})
	assert.Equal(t, search.ErrFileNameNotFound, err)
}
