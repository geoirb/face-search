package http

import (
	"github.com/valyala/fasthttp"

	search "github.com/geoirb/face-search/internal/face-search"
)

type getConfigServe struct {
	svc       search.Service
	transport *getConfigTransport
}

func (s *getConfigServe) handlerHTTP(ctx *fasthttp.RequestCtx) {
	cfg, err := s.svc.GetSearchConfig(ctx)
	s.transport.EncodeResponse(&ctx.Response, cfg, err)
}

func newGetConfigHandler(svc search.Service, transport *getConfigTransport) fasthttp.RequestHandler {
	s := getConfigServe{
		svc:       svc,
		transport: transport,
	}
	return s.handlerHTTP
}

type updateConfigServe struct {
	svc       search.Service
	transport *updateConfigTransport
}

func (s *updateConfigServe) handlerHTTP(ctx *fasthttp.RequestCtx) {
	cfg, err := s.transport.DecodeRequest(&ctx.Request)

	if err == nil {
		err = s.svc.UpdateSearchConfig(ctx, cfg)
	}
	s.transport.EncodeResponse(&ctx.Response, err)
}

func newUpdateConfigHandler(svc search.Service, transport *updateConfigTransport) fasthttp.RequestHandler {
	s := updateConfigServe{
		svc:       svc,
		transport: transport,
	}
	return s.handlerHTTP
}

type faceSearchServe struct {
	svc       search.Service
	transport *faceSearchTransport
}

func (s *faceSearchServe) handlerHTTP(ctx *fasthttp.RequestCtx) {
	sfs, err := s.transport.DecodeRequest(&ctx.Request)

	var result search.FaceSearch
	if err == nil {
		result, err = s.svc.FaceSearch(ctx, sfs)
	}
	s.transport.EncodeResponse(&ctx.Response, result, err)
}

func newFaceSearchHandler(svc search.Service, transport *faceSearchTransport) fasthttp.RequestHandler {
	s := faceSearchServe{
		svc:       svc,
		transport: transport,
	}
	return s.handlerHTTP
}

type getFaceSearchResultServe struct {
	svc       search.Service
	transport *getFaceSearchResultTransport
}

func (s *getFaceSearchResultServe) handlerHTTP(ctx *fasthttp.RequestCtx) {
	tfs, err := s.transport.DecodeRequest(ctx, &ctx.Request)

	var result search.FaceSearch
	if err == nil {
		result, err = s.svc.GetFaceSearchResult(ctx, tfs)
	}
	s.transport.EncodeResponse(&ctx.Response, result, err)
}

func newGetFaceSearchResultHandler(svc search.Service, transport *getFaceSearchResultTransport) fasthttp.RequestHandler {
	s := getFaceSearchResultServe{
		svc:       svc,
		transport: transport,
	}
	return s.handlerHTTP
}
