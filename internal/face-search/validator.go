package service

import (
	"context"
)

type Validator struct {
	svc Service
}

func NewValidator(svc Service) Service {
	return &Validator{
		svc: svc,
	}
}

func (v *Validator) GetSearchConfig(ctx context.Context) (cfg Config, err error) {
	return v.svc.GetSearchConfig(ctx)
}

func (v *Validator) UpdateSearchConfig(ctx context.Context, newCfg Config) error {
	if newCfg.Timeout == 0 || newCfg.Actions == nil {
		return errNewConfigIsNotValid
	}
	return v.svc.UpdateSearchConfig(ctx, newCfg)
}

func (v *Validator) FaceSearch(ctx context.Context, params Search) (result Result, err error) {
	if len(params.URL) == 0 {
		err = errFaceSearchParamsIsNotValid
		return
	}
	return v.svc.FaceSearch(ctx, params)
}

func (v *Validator) GetFaceSearchResult(ctx context.Context, t TaskFaceSearch) (result Result, err error) {
	return v.svc.GetFaceSearchResult(ctx, t)
}
