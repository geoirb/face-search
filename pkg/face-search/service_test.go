package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	faceSearch "github.com/geoirb/face-search/pkg/face-search"
	"github.com/geoirb/face-search/pkg/mongo"
)

var (
	testTimeout = time.Second
	testActions = []faceSearch.Action{
		{
			Type:   "navigate",
			Params: []string{"test-params-1"},
		},
		{
			Type:   "click",
			Params: []string{"test-params-2"},
		},
	}

	testSearchConfig = faceSearch.SearchConfig{
		Timeout: testTimeout,
		Actions: testActions,
	}
)

func TestGetSearchConfig(t *testing.T) {
	svc := faceSearch.NewService(
		testSearchConfig,
		time.Now().Unix,
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
	svc := faceSearch.NewService(
		faceSearch.SearchConfig{},
		time.Now().Unix,
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

// TODO test FaceSearch

func TestGetFaceSearchResult(t *testing.T) {
	// logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

	m := &mongo.Mock{}
	m.On("Get")

	// svc := faceSearch.NewService(
	// 	faceSearch.SearchConfig{},
	// 	time.Now().Unix,
	// 	testTimeout,
	// 	nil,
	// 	nil,
	// 	m,
	// 	nil,
	// 	logger,
	// )
}
