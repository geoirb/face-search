package chromedp

import (
	"fmt"
	"testing"
	"time"

	service "github.com/geoirb/face-search/internal/face-search"
	p "github.com/geoirb/face-search/internal/plugin"
	"github.com/stretchr/testify/assert"
)

var (
	testURLSelector            = "https://search4faces.com/vk01/index.html"
	testClickSelector          = "upload-button"
	testClick1Selector         = "effects-continue--upload"
	testClick2Selector         = "search-button"
	testSetUploadFilesSelector = "input[type='file']"
	testWaitNotVisibleSelector = "uppload-container"
	testWaitVisibleSelector    = "div.row.no-gutters"

	// TODO
	testFile      = "/tmp/5AYZKo5fCUk.jpg"
	testIDResult1 = "search-results1"
	testIDResult2 = "search-results2"
	testIDResult3 = "search-results3"

	testActions = []service.Action{
		{
			Type:   "navigate",
			Params: []string{testURLSelector},
		},
		{
			Type:   "click",
			Params: []string{testClickSelector},
		},
		{
			Type:   "set_upload_files",
			Params: []string{testSetUploadFilesSelector},
		},
		{
			Type:   "click",
			Params: []string{testClick1Selector},
		},
		{
			Type:   "wait_not_visible",
			Params: []string{testWaitNotVisibleSelector},
		},
		{
			Type:   "click",
			Params: []string{testClick2Selector},
		},
		{
			Type:   "wait_visible",
			Params: []string{testWaitVisibleSelector},
		},
		{
			Type:   "result_by_id",
			Params: []string{testIDResult1, testIDResult2, testIDResult3},
		},
	}

	testExpresionDir = "test-dir"
	testError        error
)

func TestActionsBuild(t *testing.T) {
	pMock := &p.Mock{}
	pMock.On("GetExpresionDir").
		Return(testExpresionDir, testError)

	c := New(
		pMock,
	)

	search := service.SearchConfig{
		Timeout:  10 * time.Minute,
		Actions:  testActions,
		FilePath: testFile,
	}
	result, err := c.Face(search)
	assert.NoError(t, err)
	fmt.Println(string(result))
}
