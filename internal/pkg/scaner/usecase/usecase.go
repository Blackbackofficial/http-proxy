package usecase

import (
	"bytes"
	"http-proxy/internal/models"
	"http-proxy/internal/pkg/proxy"
	"http-proxy/internal/pkg/scaner"
	"net/http"
)

type UseCase struct {
	repo   scaner.Repository
	hProxy proxy.HandlerProxy
}

func NewRepoUsecase(repo scaner.Repository, hProxy proxy.HandlerProxy) scaner.UseCase {
	return &UseCase{repo: repo, hProxy: hProxy}
}

func (uc *UseCase) AllRequests() ([]models.Request, models.StatusCode) {
	return uc.repo.AllRequests()
}

func (uc *UseCase) GetRequest(id int) (models.Request, models.StatusCode) {
	return uc.repo.GetRequest(id)
}

func (uc *UseCase) RepeatRequest(id int) (models.Response, models.StatusCode) {
	request, status := uc.GetRequest(id)
	if status != models.Okey {
		return models.Response{}, models.NotFound
	}

	body := bytes.NewBufferString(request.Body)
	urlStr := request.Scheme + "://" + request.Host + request.Path
	req, err := http.NewRequest(request.Method, urlStr, body)
	if err != nil {
		return models.Response{}, models.InternalError
	}

	for key, value := range request.Header {
		req.Header.Add(key, value)
	}
	resp := uc.hProxy.Proxy(req)

	return resp, models.Okey
}
