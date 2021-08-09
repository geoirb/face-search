package service_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"

	"github.com/geoirb/face-search/internal/chromedp"
	search "github.com/geoirb/face-search/internal/face-search"
	"github.com/geoirb/face-search/internal/file"
	"github.com/geoirb/face-search/internal/parser"
	"github.com/geoirb/face-search/internal/result"
)

var (
	testTimeout = time.Second
	testActions = []search.Action{
		{
			Type:   "navigate",
			Params: []string{"test-params-1"},
		},
		{
			Type:   "click",
			Params: []string{"test-params-2"},
		},
	}

	testSearchConfig = search.Config{
		Timeout: testTimeout,
		Actions: testActions,
	}

	testResult = search.Result{
		Status:    search.Success,
		UUID:      testUUID,
		PhotoHash: testFileHash,
		Profiles:  testProfiles,
		CreateAt:  testTimestamp,
		UpdateAt:  testTimestamp,
	}

	testUUID     = "test-uuid"
	testFileHash = "test-hash"
	testProfiles = []search.Profile{
		{
			FullName:    "test-name-1",
			LinkProfile: "test-link-profile-1",
			LinkPhoto:   "test-link-photo-1",
			Confidence:  "test-confidence-1",
		},
		{
			FullName:    "test-name-2",
			LinkProfile: "test-link-profile-2",
			LinkPhoto:   "test-link-photo-2",
			Confidence:  "test-confidence-2",
		},
	}
	testTimestamp int64 = 1

	testURL      = "test-url"
	testFilePath = "test-file-path"
	testPayload  = []byte("test-payload")
	errNilTest   error
	errTest      error = errors.New("test-error")
)

func TestGetSearchConfig(t *testing.T) {
	svc := search.NewService(
		testSearchConfig,
		testTimeout,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	expectedSearchConfig := testSearchConfig
	actualSearchConfig, err := svc.GetSearchConfig(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedSearchConfig, actualSearchConfig)
}

func TestUpdateSearchConfig(t *testing.T) {
	svc := search.NewService(
		search.Config{},
		testTimeout,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	err := svc.UpdateSearchConfig(context.Background(), testSearchConfig)
	assert.NoError(t, err)

	expectedSearchConfig := testSearchConfig
	actualSearchConfig, err := svc.GetSearchConfig(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedSearchConfig, actualSearchConfig)
}

func TestGetFaceSearchResult(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

		filter := search.ResultFilter{
			UUID: &testUUID,
		}

		r := &result.Mock{}
		r.On("Get", filter).
			Return(r, errNilTest)
		r.On("GetData").
			Return(testResult)

		svc := search.NewService(
			search.Config{},
			testTimeout,
			nil,
			r,
			nil,
			nil,
			logger,
		)

		expectedResult := testResult
		tfs := search.TaskFaceSearch{
			UUID: testUUID,
		}

		actualResult, err := svc.GetFaceSearchResult(context.Background(), tfs)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, actualResult)
	})

	t.Run("failed", func(t *testing.T) {
		logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

		filter := search.ResultFilter{
			UUID: &testUUID,
		}

		r := &result.Mock{}
		r.On("Get", filter).
			Return(r, errTest)

		svc := search.NewService(
			search.Config{},
			testTimeout,
			nil,
			r,
			nil,
			nil,
			logger,
		)

		tfs := search.TaskFaceSearch{
			UUID: testUUID,
		}

		_, err := svc.GetFaceSearchResult(context.Background(), tfs)
		assert.Error(t, err)
		assert.Equal(t, errTest, err)
	})
}

func TestFaceSearch(t *testing.T) {
	t.Run("new search", func(t *testing.T) {
		logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

		f0 := search.File{
			URL: testURL,
		}
		f := &file.Mock{}
		f.On("GetPath", f0).
			Return(testFilePath, errNilTest)
		f1 := search.File{
			Path: testFilePath,
		}
		f.On("GetHash", f1).
			Return(testFileHash, errNilTest)
		f.On("Delete", f1).
			Return(errNilTest)

		filter := search.ResultFilter{
			PhotoHash: &testFileHash,
		}
		r := &result.Mock{}
		r.On("Get", filter).
			Return(r, search.ErrFaceSearchResultNotFound)
		r.On("New", testFileHash).
			Return(errNilTest)
		r.On("GetStatus").
			Return("")
		r.On("SetInProgress").
			Return(errNilTest)
		r.On("GetUUID").
			Return(testUUID)
		r.On("GetData").
			Return(testResult)
		r.On("SetSuccess", testProfiles).
			Return(errNilTest)

		c := &chromedp.Mock{}
		cfg := search.Config{
			Timeout:  testSearchConfig.Timeout,
			Actions:  testSearchConfig.Actions,
			FilePath: testFilePath,
		}
		c.On("Face", cfg).
			Return(testPayload, errNilTest)

		p := &parser.Mock{}
		p.On("GetProfileList", testPayload).
			Return(testProfiles, errNilTest)

		svc := search.NewService(
			testSearchConfig,
			testTimeout,
			f,
			r,
			c,
			p,
			logger,
		)

		s := search.Search{
			File: f0,
		}

		expectedResult := testResult

		actualResult, err := svc.FaceSearch(context.Background(), s)
		time.Sleep(5 * time.Second)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, actualResult)
	})
	t.Run("priveos search is success", func(t *testing.T) {
		logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

		f0 := search.File{
			URL: testURL,
		}
		f := &file.Mock{}
		f.On("GetPath", f0).
			Return(testFilePath, errNilTest)
		f1 := search.File{
			Path: testFilePath,
		}
		f.On("GetHash", f1).
			Return(testFileHash, errNilTest)

		filter := search.ResultFilter{
			PhotoHash: &testFileHash,
		}
		r := &result.Mock{}
		r.On("Get", filter).
			Return(r, errNilTest)
		r.On("GetStatus").
			Return(search.Success)
		r.On("GetData").
			Return(testResult)

		svc := search.NewService(
			search.Config{},
			testTimeout,
			f,
			r,
			nil,
			nil,
			logger,
		)

		s := search.Search{
			File: f0,
		}

		expectedResult := testResult

		actualResult, err := svc.FaceSearch(context.Background(), s)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, actualResult)
	})
	t.Run("fail searcher.Face", func(t *testing.T) {
		logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

		f0 := search.File{
			URL: testURL,
		}
		f := &file.Mock{}
		f.On("GetPath", f0).
			Return(testFilePath, errNilTest)
		f1 := search.File{
			Path: testFilePath,
		}
		f.On("GetHash", f1).
			Return(testFileHash, errNilTest)
		f.On("Delete", f1).
			Return(errNilTest)

		filter := search.ResultFilter{
			PhotoHash: &testFileHash,
		}
		r := &result.Mock{}
		r.On("Get", filter).
			Return(r, search.ErrFaceSearchResultNotFound)
		r.On("New", testFileHash).
			Return(errNilTest)
		r.On("GetStatus").
			Return("")
		r.On("SetInProgress").
			Return(errNilTest)
		r.On("GetUUID").
			Return(testUUID)
		r.On("GetData").
			Return(testResult)
		r.On("SetFailed", errTest).
			Return(errNilTest)

		c := &chromedp.Mock{}
		cfg := search.Config{
			Timeout:  testSearchConfig.Timeout,
			Actions:  testSearchConfig.Actions,
			FilePath: testFilePath,
		}
		c.On("Face", cfg).
			Return([]byte{}, errTest)

		svc := search.NewService(
			testSearchConfig,
			testTimeout,
			f,
			r,
			c,
			nil,
			logger,
		)

		s := search.Search{
			File: f0,
		}

		expectedResult := testResult

		actualResult, err := svc.FaceSearch(context.Background(), s)
		time.Sleep(5 * time.Second)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, actualResult)
	})
	t.Run("fail parser", func(t *testing.T) {
		logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

		f0 := search.File{
			URL: testURL,
		}
		f := &file.Mock{}
		f.On("GetPath", f0).
			Return(testFilePath, errNilTest)
		f1 := search.File{
			Path: testFilePath,
		}
		f.On("GetHash", f1).
			Return(testFileHash, errNilTest)
		f.On("Delete", f1).
			Return(errNilTest)

		filter := search.ResultFilter{
			PhotoHash: &testFileHash,
		}
		r := &result.Mock{}
		r.On("Get", filter).
			Return(r, search.ErrFaceSearchResultNotFound)
		r.On("New", testFileHash).
			Return(errNilTest)
		r.On("GetStatus").
			Return("")
		r.On("SetInProgress").
			Return(errNilTest)
		r.On("GetUUID").
			Return(testUUID)
		r.On("GetData").
			Return(testResult)
		r.On("SetFailed", errTest).
			Return(errNilTest)

		c := &chromedp.Mock{}
		cfg := search.Config{
			Timeout:  testSearchConfig.Timeout,
			Actions:  testSearchConfig.Actions,
			FilePath: testFilePath,
		}
		c.On("Face", cfg).
			Return(testPayload, errNilTest)

		p := &parser.Mock{}
		p.On("GetProfileList", testPayload).
			Return([]search.Profile{}, errTest)

		svc := search.NewService(
			testSearchConfig,
			testTimeout,
			f,
			r,
			c,
			p,
			logger,
		)

		s := search.Search{
			File: f0,
		}

		expectedResult := testResult

		actualResult, err := svc.FaceSearch(context.Background(), s)
		time.Sleep(5 * time.Second)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, actualResult)
	})
}
