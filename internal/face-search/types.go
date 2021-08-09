package service

import (
	"time"
)

var (
	InProccess = "in_progress"
	Failed     = "fail"
	Success    = "success"
)

// Result ...
type Result struct {
	Status    string
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

// Config ...
type Config struct {
	Timeout  time.Duration
	Actions  []Action
	FilePath string
}

// Action for search.
type Action struct {
	Type   string
	Params []string
}

// ResultFilter ...
type ResultFilter struct {
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
