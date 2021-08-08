package service_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/geoirb/face-search/internal/chromedp"
	face "github.com/geoirb/face-search/internal/face-search"
	service "github.com/geoirb/face-search/internal/face-search"
	"github.com/geoirb/face-search/internal/file"
	"github.com/geoirb/face-search/internal/mongo"
	"github.com/geoirb/face-search/internal/parser"
)

var (
	testTimeout = time.Second
	testActions = []face.Action{
		{
			Type:   "navigate",
			Params: []string{"test-params-1"},
		},
		{
			Type:   "click",
			Params: []string{"test-params-2"},
		},
	}

	testSearchConfig = face.SearchConfig{
		Timeout: testTimeout,
		Actions: testActions,
	}

	testResult = service.Result{
		Status:    service.Success,
		UUID:      testUUID,
		PhotoHash: testFileHash,
		Profiles:  testProfiles,
		CreateAt:  testTimestamp,
		UpdateAt:  testTimestamp,
	}

	testUUID     = "test-uuid"
	testFileHash = "test-hash"
	testProfiles = []service.Profile{
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

	testUUIDFunc = func() string {
		return testUUID
	}

	testTimeFunc = func() int64 {
		return testTimestamp
	}
)

func TestGetSearchConfig(t *testing.T) {
	svc := face.NewService(
		testSearchConfig,
		time.Now().Unix,
		uuid.NewString,
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
	svc := face.NewService(
		face.SearchConfig{},
		time.Now().Unix,
		uuid.NewString,
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

		filter := service.FaceSearchFilter{
			UUID: &testUUID,
		}
		expectedResult := testResult
		m := &mongo.Mock{}
		m.On("Get", filter).
			Return(expectedResult, errNilTest)

		svc := face.NewService(
			face.SearchConfig{},
			time.Now().Unix,
			uuid.NewString,
			testTimeout,
			nil,
			nil,
			m,
			nil,
			logger,
		)

		tfs := service.TaskFaceSearch{
			UUID: testUUID,
		}

		actualResult, err := svc.GetFaceSearchResult(context.Background(), tfs)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, actualResult)
	})

	t.Run("failed", func(t *testing.T) {
		logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

		filter := service.FaceSearchFilter{
			UUID: &testUUID,
		}

		m := &mongo.Mock{}
		m.On("Get", filter).
			Return(service.Result{}, errTest)

		svc := face.NewService(
			face.SearchConfig{},
			time.Now().Unix,
			uuid.NewString,
			testTimeout,
			nil,
			nil,
			m,
			nil,
			logger,
		)

		tfs := service.TaskFaceSearch{
			UUID: testUUID,
		}

		_, err := svc.GetFaceSearchResult(context.Background(), tfs)
		assert.Error(t, err)
		assert.Equal(t, errTest, err)
	})
}

func TestFaceSearch(t *testing.T) {
	t.Run("success new face search", func(t *testing.T) {
		logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

		f0 := service.File{
			URL: testURL,
		}
		f1 := service.File{
			Path: testFilePath,
		}
		fMock := &file.Mock{}
		fMock.On("GetPath", f0).
			Return(testFilePath, errNilTest)
		fMock.On("GetHash", f1).
			Return(testFileHash, errNilTest)
		fMock.On("Delete", f1).
			Return(errNilTest)

		filter := service.FaceSearchFilter{
			PhotoHash: &testFileHash,
		}

		m := &mongo.Mock{}
		m.On("Get", filter).
			Return(service.Result{}, errNilTest)
		testResult = service.Result{
			Status:    service.Success,
			UUID:      testUUID,
			PhotoHash: testFileHash,
			Profiles:  testProfiles,
			CreateAt:  testTimestamp,
		}
		m.On("Save", testResult).
			Return(errNilTest)

		sCfg := service.SearchConfig{
			Timeout:  testTimeout,
			Actions:  testActions,
			FilePath: testFilePath,
		}

		searcherMock := &chromedp.Mock{}
		searcherMock.On("Face", sCfg).
			Return(testPayload, errNilTest)

		parser := &parser.Mock{}
		parser.On("GetProfileList", testPayload).
			Return(testProfiles, errNilTest)

		svc := face.NewService(
			testSearchConfig,
			testTimeFunc,
			testUUIDFunc,
			testTimeout,
			fMock,
			searcherMock,
			m,
			parser,
			logger,
		)

		tfs := service.Search{
			File: service.File{
				URL: testURL,
			},
		}

		expectedResult := service.Result{
			Status:    service.Fail,
			UUID:      testUUID,
			PhotoHash: testFileHash,
		}

		actualResult, err := svc.FaceSearch(context.Background(), tfs)
		time.Sleep(5 * time.Second)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, actualResult)
	})
	t.Run("success repeated face search", func(t *testing.T) {
		logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

		f0 := service.File{
			URL: testURL,
		}
		f1 := service.File{
			Path: testFilePath,
		}
		fMock := &file.Mock{}
		fMock.On("GetPath", f0).
			Return(testFilePath, errNilTest)
		fMock.On("GetHash", f1).
			Return(testFileHash, errNilTest)
		fMock.On("Delete", f1).
			Return(errNilTest)

		filter := service.FaceSearchFilter{
			PhotoHash: &testFileHash,
		}

		failedResult := service.Result{
			Status:    service.Fail,
			Error:     errTest.Error(),
			UUID:      testUUID,
			PhotoHash: testFileHash,
			CreateAt:  testTimestamp,
		}

		m := &mongo.Mock{}
		m.On("Get", filter).
			Return(failedResult, errNilTest)
		testResult = service.Result{
			Status:    service.Success,
			UUID:      testUUID,
			PhotoHash: testFileHash,
			Profiles:  testProfiles,
			CreateAt:  testTimestamp,
			UpdateAt:  testTimestamp,
		}
		m.On("Save", testResult).
			Return(errNilTest)

		sCfg := service.SearchConfig{
			Timeout:  testTimeout,
			Actions:  testActions,
			FilePath: testFilePath,
		}

		searcherMock := &chromedp.Mock{}
		searcherMock.On("Face", sCfg).
			Return(testPayload, errNilTest)

		parser := &parser.Mock{}
		parser.On("GetProfileList", testPayload).
			Return(testProfiles, errNilTest)

		svc := face.NewService(
			testSearchConfig,
			testTimeFunc,
			testUUIDFunc,
			testTimeout,
			fMock,
			searcherMock,
			m,
			parser,
			logger,
		)

		tfs := service.Search{
			File: service.File{
				URL: testURL,
			},
		}

		expectedResult := failedResult

		actualResult, err := svc.FaceSearch(context.Background(), tfs)
		time.Sleep(5 * time.Second)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, actualResult)
	})
	t.Run("failed face search", func(t *testing.T) {
		logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

		f0 := service.File{
			URL: testURL,
		}
		f1 := service.File{
			Path: testFilePath,
		}
		fMock := &file.Mock{}
		fMock.On("GetPath", f0).
			Return(testFilePath, errNilTest)
		fMock.On("GetHash", f1).
			Return(testFileHash, errNilTest)
		fMock.On("Delete", f1).
			Return(errNilTest)

		filter := service.FaceSearchFilter{
			PhotoHash: &testFileHash,
		}

		failedResult := service.Result{
			Status:    service.Fail,
			Error:     errTest.Error(),
			UUID:      testUUID,
			PhotoHash: testFileHash,
			Profiles:  nil,
			CreateAt:  testTimestamp,
		}

		m := &mongo.Mock{}
		m.On("Get", filter).
			Return(failedResult, errNilTest)
		testResult = service.Result{
			Status:    service.Fail,
			UUID:      testUUID,
			Error:     errTest.Error(),
			PhotoHash: testFileHash,
			CreateAt:  testTimestamp,
			UpdateAt:  testTimestamp,
		}
		m.On("Save", testResult).
			Return(errNilTest)

		sCfg := service.SearchConfig{
			Timeout:  testTimeout,
			Actions:  testActions,
			FilePath: testFilePath,
		}

		searcherMock := &chromedp.Mock{}
		searcherMock.On("Face", sCfg).
			Return([]byte{}, errTest)

		parser := &parser.Mock{}

		svc := face.NewService(
			testSearchConfig,
			testTimeFunc,
			testUUIDFunc,
			testTimeout,
			fMock,
			searcherMock,
			m,
			parser,
			logger,
		)

		tfs := service.Search{
			File: service.File{
				URL: testURL,
			},
		}

		expectedResult := failedResult

		actualResult, err := svc.FaceSearch(context.Background(), tfs)
		time.Sleep(5 * time.Second)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, actualResult)
	})
	t.Run("failed parse face search ", func(t *testing.T) {
		logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

		f0 := service.File{
			URL: testURL,
		}
		f1 := service.File{
			Path: testFilePath,
		}
		fMock := &file.Mock{}
		fMock.On("GetPath", f0).
			Return(testFilePath, errNilTest)
		fMock.On("GetHash", f1).
			Return(testFileHash, errNilTest)
		fMock.On("Delete", f1).
			Return(errNilTest)

		filter := service.FaceSearchFilter{
			PhotoHash: &testFileHash,
		}

		failedResult := service.Result{
			Status:    service.Fail,
			Error:     errTest.Error(),
			UUID:      testUUID,
			PhotoHash: testFileHash,
			CreateAt:  testTimestamp,
		}

		m := &mongo.Mock{}
		m.On("Get", filter).
			Return(failedResult, errNilTest)
		testResult = service.Result{
			Status:    service.Fail,
			UUID:      testUUID,
			Error:     errTest.Error(),
			PhotoHash: testFileHash,
			Profiles:  nil,
			CreateAt:  testTimestamp,
			UpdateAt:  testTimestamp,
		}
		m.On("Save", testResult).
			Return(errNilTest)

		sCfg := service.SearchConfig{
			Timeout:  testTimeout,
			Actions:  testActions,
			FilePath: testFilePath,
		}

		searcherMock := &chromedp.Mock{}
		searcherMock.On("Face", sCfg).
			Return(testPayload, errNilTest)

		parser := &parser.Mock{}
		parser.On("GetProfileList", testPayload).
			Return([]service.Profile{}, errTest)

		svc := face.NewService(
			testSearchConfig,
			testTimeFunc,
			testUUIDFunc,
			testTimeout,
			fMock,
			searcherMock,
			m,
			parser,
			logger,
		)

		tfs := service.Search{
			File: service.File{
				URL: testURL,
			},
		}

		expectedResult := failedResult

		actualResult, err := svc.FaceSearch(context.Background(), tfs)
		time.Sleep(5 * time.Second)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, actualResult)
	})
}
