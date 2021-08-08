package http

import (
	"net/http"

	"github.com/fasthttp/router"

	service "github.com/geoirb/face-search/internal/face-search"
)

const (
	version = "/v1/api"

	getConfigURI           = version + "/config"
	updateConfigURI        = version + "/config"
	faceSearchURI          = version + "/face_search"
	getFaceSearchResultURI = version + "/face_search/{uuid}"
)

type builder func(payload interface{}, err error) ([]byte, error)

// Routing adds handles to router.
func Routing(router *router.Router, svc service.Service, builder builder) {
	router.Handle(http.MethodGet, getConfigURI, newGetConfigHandler(svc, newGetConfigTransport(builder)))
	router.Handle(http.MethodPut, updateConfigURI, newUpdateConfigHandler(svc, newUpdateConfigTransport(builder)))
	router.Handle(http.MethodPost, faceSearchURI, newFaceSearchHandler(svc, newFaceSearchTransport(builder)))
	router.Handle(http.MethodGet, getFaceSearchResultURI, newGetFaceSearchResultHandler(svc, newGetFaceSearchResultTransport(builder)))
}
