package http

import (
	"time"
)

type searchConfig struct {
	Timeout time.Duration `json:"timeout"`
	Actions []action      `json:"actions"`
}

type action struct {
	Type   string   `json:"type"`
	Params []string `json:"params"`
}

type startFaceSearch struct {
	URL string `json:"url"`
}

type faceSearch struct {
	Status    string    `json:"status"`
	Error     string    `json:"error,omitempty"`
	UUID      string    `json:"uuid"`
	PhotoHash string    `json:"photo_hash"`
	Profiles  []profile `json:"profiles,omitempty"`
	UpdateAt  int64     `json:"update_at,omitempty"`
	CreateAt  int64     `json:"create_at,omitempty"`
}

type profile struct {
	FullName    string `json:"full_name"`
	LinkProfile string `json:"link_profile"`
	LinkPhoto   string `json:"link_photo"`
	Confidence  string `json:"confidence"`
}
