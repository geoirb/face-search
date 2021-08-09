package service

import "errors"

// errors
var (
	ErrResult                     = errors.New("result not created")
	ErrFaceSearchResultNotFound   = errors.New("face search result not found")
	ErrFileNotFound               = errors.New("file not found")
	ErrFileNameNotFound           = errors.New("file name not found")
	errFaceSearchParamsIsNotValid = errors.New("face search params is not valid")
)
