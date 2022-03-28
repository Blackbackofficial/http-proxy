package proxy

import (
	"http-proxy/internal/models"
	"net/http"
)

type RepoProxy interface {
	SaveRequest(r *http.Request) (int, error)
	SaveResponse(reqId int, resp *http.Response) (models.Response, error)
}
type HandlerProxy interface {
	ProxyHTTP(r *http.Request) models.Response
}
