package service

import (
	"context"
	"fmt"
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
	isValid := false
	errStr := "config is not valid: "
	if newCfg.Timeout == 0 {
		errStr += "timeout not found "
		isValid = false
	}
	if len(newCfg.Actions) == 0 {
		errStr += "actions not found "
		isValid = false
	}
	for i, action := range newCfg.Actions {
		if len(action.Type) == 0 {
			errStr += fmt.Sprintf(" actions %d: type unknown", i)
			isValid = false
		}
		if len(action.Params) == 0 {
			errStr += fmt.Sprintf(" actions %d: params empty", i)
			isValid = false
		}
	}
	if isValid {
		return v.svc.UpdateSearchConfig(ctx, newCfg)
	}
	return fmt.Errorf(errStr)
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
