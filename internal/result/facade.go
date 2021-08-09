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

// Facade result.
type Facade struct {
	search.Result

	timeFunc func() int64
	uuidFunc func() string

	storage storage
}

// NewFacade ...
func NewFacade(
	timeFunc func() int64,
	uuidFunc func() string,

	storage storage,
) search.IResult {
	return &Facade{
		timeFunc: timeFunc,
		uuidFunc: uuidFunc,

		storage: storage,
	}
}

// Get result by filter
func (p *Facade) Get(ctx context.Context, filter search.ResultFilter) (search.IResult, error) {
	r := &Facade{
		timeFunc: p.timeFunc,
		uuidFunc: p.uuidFunc,
		storage:  p.storage,
	}
	err := r.get(ctx, filter)
	return r, err
}

// New data of result.
func (r *Facade) New(ctx context.Context, hash string) error {
	r.UUID = r.uuidFunc()
	r.PhotoHash = hash
	r.CreateAt = r.timeFunc()
	return r.storage.Save(ctx, r.Result)
}

// GetStatus of result data.
func (r *Facade) GetStatus() string {
	return r.Status
}

// GetUUID of result data.
func (r *Facade) GetUUID() string {
	return r.UUID
}

// GetData of result data.
func (r *Facade) GetData() search.Result {
	return r.Result
}

// SetInProgress status of result.
func (r *Facade) SetInProgress(ctx context.Context) error {
	r.Status = search.InProccess
	r.Error = ""
	r.Profiles = nil
	r.UpdateAt = r.timeFunc()
	return r.storage.Update(ctx, r.Result)
}

// SetSuccess status of result.
func (r *Facade) SetSuccess(ctx context.Context, profiles []search.Profile) error {
	r.Status = search.Success
	r.Error = ""
	r.Profiles = profiles
	r.UpdateAt = r.timeFunc()
	return r.storage.Update(ctx, r.Result)
}

// SetFailed status of result.
func (r *Facade) SetFailed(ctx context.Context, err error) error {
	r.Status = search.Failed
	r.Error = err.Error()
	r.Profiles = nil
	r.UpdateAt = r.timeFunc()
	return r.storage.Update(ctx, r.Result)
}

func (r *Facade) get(ctx context.Context, filter search.ResultFilter) (err error) {
	r.Result, err = r.storage.Get(ctx, filter)
	return
}
