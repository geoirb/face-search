package plugin

import (
	_ "embed"
	"fmt"
	"os"
)

var (
	//go:embed manifest.json
	manifestJSON string
	//go:embed background.js
	backgroundJS string
)

type proxy interface {
	Get() (Proxy, error)
}

type Plugin struct {
	proxy proxy

	pluginDirLayout string
}

func New(
	proxy proxy,
	pluginDirLayout string,
) *Plugin {
	return &Plugin{
		proxy:           proxy,
		pluginDirLayout: pluginDirLayout,
	}
}

func (p *Plugin) GetExpresionDir() (expresionDir string, err error) {
	proxy, err := p.proxy.Get()
	if err != nil {
		return
	}
	expresionDir = fmt.Sprintf(p.pluginDirLayout, proxy.ID)
	backgroundJS = fmt.Sprintf(backgroundJS, proxy.IP, proxy.Port, proxy.Login, proxy.Password)

	if _, err = os.Stat(expresionDir); os.IsNotExist(err) {
		if err = os.Mkdir(expresionDir, 0750); err != nil {
			return
		}

		var file *os.File
		if file, err = os.Create(expresionDir + "/manifest.json"); err != nil {
			return
		}
		if _, err = file.WriteString(manifestJSON); err != nil {
			return
		}
		if err = file.Close(); err != nil {
			return
		}

		if file, err = os.Create(expresionDir + "/background.js"); err != nil {
			return
		}
		if _, err = file.WriteString(backgroundJS); err != nil {
			return
		}
		if err = file.Close(); err != nil {
			return
		}
	}

	return
}
