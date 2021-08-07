package chromedp

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"gopkg.in/yaml.v2"

	service "github.com/geoirb/face-search/pkg/face-search"
)

type buildFunc func(params []string) (chromedp.Action, error)

var actionFunc = map[string]buildFunc{
	"navigate":         navigate,
	"click":            click,
	"wait_not_visible": waitNotVisible,
	"wait_visible":     waitVisible,
	"sleep":             sleep,
}

func (c *Chromedp) actionsBuild(searchActions []service.Action, file string, result interface{}) (actions []chromedp.Action, err error) {
	actions = make([]chromedp.Action, 0, len(actions))
	for _, action := range searchActions {
		var a chromedp.Action
		switch action.Type {
		case "result_by_id":
			a = resultByID(action.Params, result)
		case "set_upload_files":
			a, err = setUploadFiles(action.Params, file)
		default:
			actionFunc, ok := actionFunc[action.Type]
			if !ok {
				err = errUnknownActionType
				return
			}
			a, err = actionFunc(action.Params)
		}
		if err != nil {
			return
		}
		actions = append(actions, a)
	}
	return
}

func navigate(params []string) (a chromedp.Action, err error) {
	if len(params) != 1 {
		err = errors.New("navigate: wrong number of params")
		return
	}
	a = chromedp.Navigate(params[0])
	return
}

func click(params []string) (a chromedp.Action, err error) {
	if len(params) != 1 {
		err = errors.New("click: wrong number of params")
		return
	}
	a = chromedp.Click(params[0], chromedp.NodeEnabled, chromedp.NodeVisible, chromedp.BySearch)
	return
}

func setUploadFiles(params []string, file string) (a chromedp.Action, err error) {
	if len(params) != 1 {
		err = errors.New("set_upload_files: wrong number of params")
		return
	}
	a = chromedp.SetUploadFiles(params[0], []string{file}, chromedp.NodeVisible, chromedp.BySearch)
	return
}

func waitNotVisible(params []string) (a chromedp.Action, err error) {
	if len(params) != 1 {
		err = errors.New("wait_not_visible: wrong number of params")
		return
	}
	a = chromedp.WaitNotVisible(params[0], chromedp.BySearch)
	return
}

func waitVisible(params []string) (a chromedp.Action, err error) {
	if len(params) != 1 {
		err = errors.New("wait_visible: wrong number of params")
		return
	}
	a = chromedp.WaitVisible(params[0], chromedp.BySearch)
	return
}

func sleep(params []string) (a chromedp.Action, err error) {
	if len(params) != 1 {
		err = errors.New("wait: wrong number of params")
		return
	}
	var timeout time.Duration
	err = yaml.Unmarshal([]byte(params[0]), &timeout)
	a = chromedp.Sleep(timeout)
	return
}

func resultByID(params []string, result interface{}) (a chromedp.Action) {
	a = chromedp.EvaluateAsDevTools(getElementByIDs(params), result)
	return
}

func getElementByIDs(ids []string) string {
	output := make([]string, 0, len(ids))
	for _, id := range ids {
		output = append(
			output,
			fmt.Sprintf(`document.getElementById("%s").innerHTML`, id),
		)
	}
	return strings.Join(output, "+")
}
