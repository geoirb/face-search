package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Service interface {
	// GetSearchConfig returns current search config.
	GetSearchConfig(ctx context.Context) (cfg SearchConfig, err error)
	// UpdateSearchConfig ...
	UpdateSearchConfig(ctx context.Context, newSearch SearchConfig) error
	// FaceSearch start face search, if the previous search for params failed
	// or params is new.
	// If the previous search for params success, FaceSearch returns result of face search.
	FaceSearch(ctx context.Context, params Search) (result Result, err error)
	// GetFaceSearchResult returns results face search by tfs.
	GetFaceSearchResult(ctx context.Context, tfs TaskFaceSearch) (result Result, err error)
}

type file interface {
	GetPath(file File) (path string, err error)
	GetHash(file File) (hash string, err error)
	Delete(file File) (err error)
}

type storage interface {
	Save(ctx context.Context, result Result) error
	Get(ctx context.Context, filter FaceSearchFilter) (Result, error)
}

type searcher interface {
	Face(search SearchConfig) (result []byte, err error)
}

type parser interface {
	GetProfileList(payload []byte) ([]Profile, error)
}

type service struct {
	search             SearchConfig
	timeFunc           func() int64
	uuidFunc           func() string
	storageSaveTimeout time.Duration

	file     file
	searcher searcher
	storage  storage
	parser   parser

	queue chan struct{}

	logger log.Logger
}

// NewService returns face search service.
func NewService(
	searchConfig SearchConfig,
	timeFunc func() int64,
	uuidFunc func() string,
	storageSaveTimeout time.Duration,

	file file,
	searcher searcher,
	storage storage,
	parser parser,

	logger log.Logger,
) Service {
	return &service{
		search:             searchConfig,
		timeFunc:           timeFunc,
		uuidFunc:           uuidFunc,
		storageSaveTimeout: storageSaveTimeout,

		file:     file,
		searcher: searcher,
		storage:  storage,
		parser:   parser,

		queue: make(chan struct{}, 1),

		logger: logger,
	}
}

func (s *service) GetSearchConfig(ctx context.Context) (SearchConfig, error) {
	return s.search, nil
}

func (s *service) UpdateSearchConfig(ctx context.Context, cfg SearchConfig) error {
	if cfg.Timeout != 0 {
		s.search.Timeout = cfg.Timeout
	}
	if len(cfg.Actions) != 0 {
		s.search.Actions = cfg.Actions
	}
	return nil
}

func (s *service) FaceSearch(ctx context.Context, sfs Search) (result Result, err error) {
	logger := log.WithPrefix(s.logger, "method", "FaceSearch")

	path, err := s.file.GetPath(sfs.File)
	if err != nil {
		level.Error(logger).Log("get file path", err)
		return
	}

	file := File{
		Path: path,
	}
	hash, err := s.file.GetHash(file)
	if err != nil {
		level.Error(logger).Log("get file hash", "path", path, "err", err)
		return
	}

	filter := FaceSearchFilter{
		PhotoHash: &hash,
	}
	if result, err = s.storage.Get(ctx, filter); err == nil && result.Status == Success {
		return
	}
	if err != nil && err != ErrFaceSearchResultNotFound {
		level.Error(logger).Log("get result from db", "hash", hash, "err", err)
		return
	}
	err = nil

	result.Status = InProccess
	result.PhotoHash = hash
	if result.CreateAt == 0 {
		result.CreateAt = s.timeFunc()
	} else {
		result.UpdateAt = s.timeFunc()
	}

	if len(result.UUID) == 0 {
		result.UUID = s.uuidFunc()
	}
	if err = s.storage.Save(ctx, result); err != nil {
		level.Error(logger).Log("msg", "save result of face search", "err", err)
	}

	go func(result Result, file File, logger log.Logger) {
		s.queue <- struct{}{}
		defer func() {
			s.file.Delete(file)
			<-s.queue
			level.Info(logger).Log("msg", "finish face search")
		}()
		level.Info(logger).Log("msg", "start face search")

		search := SearchConfig{
			Timeout:  s.search.Timeout,
			Actions:  s.search.Actions,
			FilePath: file.Path,
		}
		payloadResult, err := s.searcher.Face(search)
		if err != nil {
			result.Status = Fail
			result.Error = err.Error()
			level.Error(logger).Log("msg", "face search", "err", err)
		} else {
			result.Status = Success
			result.Error = ""
			if result.Profiles, err = s.parser.GetProfileList(payloadResult); err != nil {
				result.Status = Fail
				result.Error = err.Error()
				result.Profiles = nil
				level.Error(logger).Log("msg", "parse face search result", "err", err)
			}
		}

		if result.CreateAt == 0 {
			result.CreateAt = s.timeFunc()
		} else {
			result.UpdateAt = s.timeFunc()
		}

		ctx, cancel := context.WithTimeout(context.Background(), s.storageSaveTimeout)
		defer cancel()
		if err = s.storage.Save(ctx, result); err != nil {
			level.Error(logger).Log("msg", "save result of face search", "err", err)
		}
	}(result, file, log.WithPrefix(logger, "uuid", result.UUID))
	return
}

func (s *service) GetFaceSearchResult(ctx context.Context, tsf TaskFaceSearch) (r Result, err error) {
	logger := log.WithPrefix(s.logger, "method", "GetFaceSearchResult")

	filter := FaceSearchFilter{
		UUID: &tsf.UUID,
	}
	if r, err = s.storage.Get(ctx, filter); err != nil {
		level.Error(logger).Log("get result from db", "uuid", filter.UUID, "err", err)
	}
	return
}
