package service

import "errors"

// errors
var (
	ErrFaceSearchResultNotFound = errors.New("face search result not found")
	ErrFileNotFound             = errors.New("file not found")
	ErrFileNameNotFound         = errors.New("file name not found")
)
