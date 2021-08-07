package chromedp

import (
	"github.com/chromedp/chromedp"
)

var (
	UserAgents = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.128 Safari/537.36",
	}

	chromedpOpts = append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("user-agent", UserAgents[0]),
	)
)

func (c *Chromedp) optionsBuild() ([]chromedp.ExecAllocatorOption, error) {
	dir, err := c.plugin.GetExpresionDir()
	if err != nil {
		return nil, err
	}
	
	chromedpOpts = append(chromedpOpts,
		chromedp.Flag("load-extension", dir),
	)

	return chromedpOpts, nil
}
