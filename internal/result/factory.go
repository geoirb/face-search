package result

import (
	"context"

	search "github.com/geoirb/face-search/internal/face-search"
)

type storage interface {
	Save(ctx context.Context, result search.Result) error
	Update(ctx context.Context, result search.Result) (err error)
	Get(ctx context.Context, filter search.ResultFilter) (search.Result, error)
}

type factory struct {
	timeFunc func() int64
	uuidFunc func() string

	storage storage
}

func NewFactoryFunc(
	timeFunc func() int64,
	uuidFunc func() string,

	storage storage,
) func(context.Context, search.ResultFilter) (*Facade, error) {
	factory := &factory{
		timeFunc: timeFunc,
		uuidFunc: uuidFunc,

		storage: storage,
	}
	return factory.Create
}

// Create result by filter
func (p *factory) Create(ctx context.Context, filter search.ResultFilter) (r *Facade, err error) {
	r = &Facade{
		timeFunc: p.timeFunc,
		uuidFunc: p.uuidFunc,
		storage:  p.storage,
	}
	err = r.get(ctx, filter)
	return
}
