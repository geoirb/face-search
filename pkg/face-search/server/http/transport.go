package http

import (
	"encoding/json"

	"github.com/valyala/fasthttp"

	search "github.com/geoirb/face-search/pkg/face-search"
)

type getConfigTransport struct {
	builder builder
}

func newGetConfigTransport(builder builder) *getConfigTransport {
	return &getConfigTransport{
		builder: builder,
	}
}

func (t *getConfigTransport) DecodeRequest(req *fasthttp.Request) (err error) {
	return
}

func (t *getConfigTransport) EncodeResponse(res *fasthttp.Response, sc search.SearchConfig, err error) {
	response := searchConfig{
		Timeout: sc.Timeout,
		Actions: make([]action, 0, len(sc.Actions)),
	}

	for _, a := range sc.Actions {
		response.Actions = append(response.Actions, action(a))
	}
	body, _ := t.builder(response, err)
	res.SetBody(body)
}

type updateConfigTransport struct {
	builder builder
}

func newUpdateConfigTransport(builder builder) *updateConfigTransport {
	return &updateConfigTransport{
		builder: builder,
	}
}

func (t *updateConfigTransport) DecodeRequest(req *fasthttp.Request) (cfg search.SearchConfig, err error) {
	request := searchConfig{}
	if err = json.Unmarshal(req.Body(), &request); err != nil {
		return
	}

	cfg = search.SearchConfig{
		Timeout: request.Timeout,
		Actions: make([]search.Action, 0, len(request.Actions)),
	}

	for _, a := range request.Actions {
		cfg.Actions = append(cfg.Actions, search.Action(a))
	}
	return
}

func (t *updateConfigTransport) EncodeResponse(res *fasthttp.Response, err error) {
	body, _ := t.builder(nil, err)
	res.SetBody(body)
}

type faceSearchTransport struct {
	builder builder
}

func newFaceSearchTransport(builder builder) *faceSearchTransport {
	return &faceSearchTransport{
		builder: builder,
	}
}

func (t *faceSearchTransport) DecodeRequest(req *fasthttp.Request) (sfs search.StartFaceSearch, err error) {
	request := startFaceSearch{}
	if err = json.Unmarshal(req.Body(), &request); err != nil {
		return
	}

	sfs.File.URL = request.URL
	return
}

func (t *faceSearchTransport) EncodeResponse(res *fasthttp.Response, result search.FaceSearch, err error) {
	response := faceSearch{
		Status:    result.Status,
		Error:     result.Error,
		UUID:      result.UUID,
		PhotoHash: result.PhotoHash,
		Profiles:  make([]profile, 0, len(result.Profiles)),
		UpdateAt:  result.UpdateAt,
		CreateAt:  result.CreateAt,
	}
	for _, p := range result.Profiles {
		response.Profiles = append(response.Profiles, profile(p))
	}
	body, _ := t.builder(response, err)
	res.SetBody(body)
}

type getFaceSearchResultTransport struct {
	builder builder
}

func newGetFaceSearchResultTransport(builder builder) *getFaceSearchResultTransport {
	return &getFaceSearchResultTransport{
		builder: builder,
	}
}

func (t *getFaceSearchResultTransport) DecodeRequest(ctx *fasthttp.RequestCtx, req *fasthttp.Request) (tfs search.TaskFaceSearch, err error) {
	tfs.UUID = ctx.UserValue("uuid").(string)
	return
}

func (t *getFaceSearchResultTransport) EncodeResponse(res *fasthttp.Response, result search.FaceSearch, err error) {
	response := faceSearch{
		Status:    result.Status,
		Error:     result.Error,
		UUID:      result.UUID,
		PhotoHash: result.PhotoHash,
		Profiles:  make([]profile, 0, len(result.Profiles)),
		UpdateAt:  result.UpdateAt,
		CreateAt:  result.CreateAt,
	}
	for _, p := range result.Profiles {
		response.Profiles = append(response.Profiles, profile(p))
	}
	body, _ := t.builder(response, err)
	res.SetBody(body)
}
