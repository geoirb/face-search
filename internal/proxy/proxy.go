// Todo
package proxy

import (
	"math/rand"

	"github.com/geoirb/face-search/internal/plugin"
)

var (
	proxys = []plugin.Proxy{
		{
			ID:       "first",
			IP:       "91.188.241.198",
			Port:     "9982",
			Login:    "gK1Qbz",
			Password: "EWAtCM",
		},
		{
			ID:       "second",
			IP:       "91.188.242.95",
			Port:     "9733",
			Login:    "gK1Qbz",
			Password: "EWAtCM",
		},
		{
			ID:       "third",
			IP:       "91.188.243.6",
			Port:     "9233",
			Login:    "gK1Qbz",
			Password: "EWAtCM",
		},
	}
)

// Proxy ...
type Proxy struct {
}

// New proxy.
func New() *Proxy {
	return &Proxy{}
}

// Get proxy.
func (p *Proxy) Get() (plugin.Proxy, error) {
	i := rand.Intn(len(proxys))
	return proxys[i], nil
}
