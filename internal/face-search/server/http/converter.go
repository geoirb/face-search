package http

import (
	search "github.com/geoirb/face-search/internal/face-search"
)

func toResultResponse(src search.Result) resultResponse {
	dst := resultResponse{
		Status:    src.Status,
		Error:     src.Error,
		UUID:      src.UUID,
		PhotoHash: src.PhotoHash,
		Profiles:  make([]profile, 0, len(src.Profiles)),
		UpdateAt:  src.UpdateAt,
		CreateAt:  src.CreateAt,
	}
	for _, p := range src.Profiles {
		dst.Profiles = append(dst.Profiles, profile(p))
	}
	return dst
}
