package result

import (
	"context"
	"errors"
	"testing"

	search "github.com/geoirb/face-search/internal/face-search"
	"github.com/geoirb/face-search/internal/mongo"
	"github.com/stretchr/testify/assert"
)

var (
	timeFunc = func() int64 {
		return testTimestamp
	}

	testTimestamp = int64(1)

	uuidFunc = func() string {
		return testUUID
	}

	testUUID = "test-uuid"

	testResult = search.Result{
		Status:    search.Success,
		UUID:      testUUID,
		PhotoHash: testFileHash,
		Profiles:  testProfiles,
		CreateAt:  testTimestamp,
		UpdateAt:  testTimestamp,
	}
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

	errNilTest error
	errTest    = errors.New("test-error")
)

func TestGet(t *testing.T) {
	filter := search.ResultFilter{
		UUID: &testUUID,
	}
	s := &mongo.Mock{}
	s.On("Get", filter).
		Return(testResult, errNilTest)
	r := NewFacade(
		timeFunc,
		uuidFunc,
		s,
	)

	expectedResult := &Facade{
		Result: testResult,

		timeFunc: timeFunc,
		uuidFunc: uuidFunc,
		storage:  s,
	}

	actualResult, err := r.Get(context.Background(), filter)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult.GetData(), actualResult.GetData())
}

func TestNew(t *testing.T) {
	result := search.Result{
		UUID:      testUUID,
		PhotoHash: testFileHash,
		CreateAt:  testTimestamp,
	}
	s := &mongo.Mock{}
	s.On("Save", result).
		Return(errNilTest)
	r := NewFacade(
		timeFunc,
		uuidFunc,
		s,
	)

	expectedResult := &Facade{
		Result: result,

		timeFunc: timeFunc,
		uuidFunc: uuidFunc,
		storage:  s,
	}

	err := r.New(context.Background(), testFileHash)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult.GetData(), r.GetData())
}

func TestSetInProgress(t *testing.T) {
	result := search.Result{
		Status:   search.InProccess,
		Error:    "",
		UpdateAt: testTimestamp,
	}
	s := &mongo.Mock{}
	s.On("Update", result).
		Return(errNilTest)
	r := NewFacade(
		timeFunc,
		uuidFunc,
		s,
	)

	expectedResult := &Facade{
		Result: result,

		timeFunc: timeFunc,
		uuidFunc: uuidFunc,
		storage:  s,
	}

	err := r.SetInProgress(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedResult.GetData(), r.GetData())
}

func TestSetSuccess(t *testing.T) {
	result := search.Result{
		Status:   search.Success,
		Error:    "",
		Profiles: testProfiles,
		UpdateAt: testTimestamp,
	}
	s := &mongo.Mock{}
	s.On("Update", result).
		Return(errNilTest)
	r := NewFacade(
		timeFunc,
		uuidFunc,
		s,
	)

	expectedResult := &Facade{
		Result: result,

		timeFunc: timeFunc,
		uuidFunc: uuidFunc,
		storage:  s,
	}

	err := r.SetSuccess(context.Background(), testProfiles)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult.GetData(), r.GetData())
}

func TestSetFailed(t *testing.T) {
	result := search.Result{
		Status:   search.Failed,
		Error:    errTest.Error(),
		UpdateAt: testTimestamp,
	}
	s := &mongo.Mock{}
	s.On("Update", result).
		Return(errNilTest)
	r := NewFacade(
		timeFunc,
		uuidFunc,
		s,
	)

	expectedResult := &Facade{
		Result: result,

		timeFunc: timeFunc,
		uuidFunc: uuidFunc,
		storage:  s,
	}

	err := r.SetFailed(context.Background(), errTest)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult.GetData(), r.GetData())
}
