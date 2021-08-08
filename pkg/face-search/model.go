package service

import (
	"time"
)

// FaceSearch ...
type FaceSearch struct {
	Status    bool
	Error     string
	UUID      string
	PhotoHash string
	Profiles  []Profile
	UpdateAt  int64
	CreateAt  int64
}

// Profile ...
type Profile struct {
	FullName    string
	LinkProfile string
	LinkPhoto   string
	Confidence  string
}

// SearchConfig ...
type SearchConfig struct {
	Timeout time.Duration
	Actions []Action
	FilePath    string
}

// Action for search.
type Action struct {
	Type   string
	Params []string
}

// FaceSearchFilter ...
type FaceSearchFilter struct {
	UUID      *string
	PhotoHash *string
}

// Search ...
type Search struct {
	File
}

type File struct {
	Path string
	URL  string
}

// TaskFaceSearch ...
type TaskFaceSearch struct {
	UUID string
}
