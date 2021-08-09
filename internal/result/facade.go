package result

import (
	"context"

	search "github.com/geoirb/face-search/internal/face-search"
)

// Facade result.
type Facade struct {
	search.Result

	timeFunc func() int64
	uuidFunc func() string

	storage storage
}

// New data of result.
func (r *Facade) New(ctx context.Context) error {
	r.CreateAt = r.timeFunc()
	return r.storage.Save(ctx, r.Result)
}

// SetInProgress status of result.
func (r *Facade) SetInProgress(ctx context.Context) error {
	r.Status = search.InProccess
	r.Error = ""
	r.Profiles = nil
	r.UpdateAt = r.timeFunc()
	return r.storage.Update(ctx, r.Result)
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
