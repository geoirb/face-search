package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Service interface {
	// GetSearchConfig returns current search config.
	GetSearchConfig(ctx context.Context) (cfg Config, err error)
	// UpdateSearchConfig ...
	UpdateSearchConfig(ctx context.Context, newSearch Config) error
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

type resultFactory func(context.Context, ResultFilter) (result, error)

type result interface {
	New(ctx context.Context) error
	GetStatus() string
	GetUUID() string
	GetData() Result
	SetInProgress(ctx context.Context) error
	SetSuccess(ctx context.Context, profiles []Profile) error
	SetFailed(ctx context.Context, err error) error
}

type searcher interface {
	Face(search Config) (result []byte, err error)
}

type parser interface {
	GetProfileList(payload []byte) ([]Profile, error)
}

type service struct {
	search              Config
	resultUpdateTimeout time.Duration

	file      file
	getResult resultFactory
	searcher  searcher
	parser    parser

	queue chan struct{}

	logger log.Logger
}

// NewService returns face search service.
func NewService(
	searchConfig Config,
	resultUpdateTimeout time.Duration,

	file file,
	getResult resultFactory,
	searcher searcher,
	parser parser,

	logger log.Logger,
) Service {
	return &service{
		search:              searchConfig,
		resultUpdateTimeout: resultUpdateTimeout,

		file:      file,
		getResult: getResult,
		searcher:  searcher,
		parser:    parser,

		queue: make(chan struct{}, 1),

		logger: logger,
	}
}

func (s *service) GetSearchConfig(ctx context.Context) (Config, error) {
	return s.search, nil
}

func (s *service) UpdateSearchConfig(ctx context.Context, cfg Config) error {
	if cfg.Timeout != 0 {
		s.search.Timeout = cfg.Timeout
	}
	if len(cfg.Actions) != 0 {
		s.search.Actions = cfg.Actions
	}
	return nil
}

func (s *service) FaceSearch(ctx context.Context, sh Search) (r Result, err error) {
	logger := log.WithPrefix(s.logger, "method", "FaceSearch")

	path, err := s.file.GetPath(sh.File)
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

	filter := ResultFilter{
		PhotoHash: &hash,
	}
	result, err := s.getResult(ctx, filter)
	if err == ErrFaceSearchResultNotFound {
		err = result.New(ctx)
	}
	if err != nil {
		level.Error(logger).Log("get face search result from db", "hash", hash, "err", err)
		return
	}
	if result.GetStatus() == Success {
		r = result.GetData()
		return
	}

	result.SetInProgress(ctx)
	go s.start(result, file, log.WithPrefix(logger, "uuid", result.GetUUID()))
	return
}

func (s *service) start(result result, file File, logger log.Logger) {
	s.queue <- struct{}{}
	defer func() {
		s.file.Delete(file)
		<-s.queue
		level.Info(logger).Log("msg", "finish face search")
	}()
	level.Info(logger).Log("msg", "start face search")

	search := Config{
		Timeout:  s.search.Timeout,
		Actions:  s.search.Actions,
		FilePath: file.Path,
	}
	payload, err := s.searcher.Face(search)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.resultUpdateTimeout)
		defer cancel()
		if err := result.SetFailed(ctx, err); err != nil {
			level.Error(logger).Log("msg", "result status set fail", "err", err)
		}
		return
	}
	profiles, err := s.parser.GetProfileList(payload)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.resultUpdateTimeout)
		defer cancel()
		if err := result.SetFailed(ctx, err); err != nil {
			level.Error(logger).Log("msg", "result status set fail", "err", err)
		}
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.resultUpdateTimeout)
	defer cancel()
	if err := result.SetSuccess(ctx, profiles); err != nil {
		level.Error(logger).Log("msg", "result status set success", "err", err)
	}
}

func (s *service) GetFaceSearchResult(ctx context.Context, t TaskFaceSearch) (r Result, err error) {
	logger := log.WithPrefix(s.logger, "method", "GetFaceSearchResult")

	filter := ResultFilter{
		UUID: &t.UUID,
	}
	result, err := s.getResult(ctx, filter)
	if err != nil {
		level.Error(logger).Log("get result from db", "uuid", t.UUID, "err", err)
	}
	r = result.GetData()
	return
}
