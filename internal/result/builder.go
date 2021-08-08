package result

import (
	faceSearch "github.com/geoirb/face-search/internal/face-search"
)

type Builder struct {
	timeFunc func() int64
	uuidFunc func() string
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) NewResult(hash string) faceSearch.Result {
	return faceSearch.Result{
		Status:   faceSearch.InProccess,
		UUID:     b.uuidFunc(),
		CreateAt: b.timeFunc(),
	}
}
