package chromedp

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	faceSearch "github.com/geoirb/face-search/internal/face-search"
)

type plugin interface {
	GetExpresionDir() (string, error)
}

type Chromedp struct {
	plugin plugin
}

func New(plugin plugin) *Chromedp {
	return &Chromedp{
		plugin: plugin,
	}
}

func (c *Chromedp) Face(search faceSearch.Config) (result []byte, err error) {
	timeoutContext, cancelAllocCtx, cancelNewContext, cancelWithTimeout, err := c.getTimeoutContext(search.Timeout)
	if err != nil {
		return
	}
	defer cancelAllocCtx()
	defer cancelNewContext()
	defer cancelWithTimeout()

	actions, err := c.actionsBuild(search.Actions, search.FilePath, &result)
	actions = append(actions, chromedp.Sleep(5*time.Second))
	if err != nil {
		return
	}

	err = chromedp.Run(timeoutContext, actions...)
	return
}

func (c *Chromedp) getTimeoutContext(timeout time.Duration) (context.Context, context.CancelFunc, context.CancelFunc, context.CancelFunc, error) {
	opts, err := c.optionsBuild()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	allocCtx, cancelAllocCtx := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancelNewContext := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	timeoutContext, cancelWithTimeout := context.WithTimeout(ctx, timeout)

	return timeoutContext, cancelAllocCtx, cancelNewContext, cancelWithTimeout, nil
}
