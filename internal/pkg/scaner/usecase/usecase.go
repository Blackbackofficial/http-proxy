package usecase

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"http-proxy/internal/models"
	"http-proxy/internal/pkg/proxy"
	"http-proxy/internal/pkg/scaner"
	"http-proxy/internal/pkg/utils"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
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

func (uc *UseCase) Scan(id int) models.StatusCode {
	request, status := uc.GetRequest(id)
	if status != models.Okey {
		return models.NotFound
	}
	fmt.Println(request)

	if _, err := os.Stat("./params"); errors.Is(err, os.ErrNotExist) {
		_, err := grab.Get("./params", "https://raw.githubusercontent.com/PortSwigger/param-miner/master/resources/params")
		if err != nil {
			log.Fatal(err)
		}
	}

	body := bytes.NewBufferString(request.Body)
	urlStr := request.Scheme + "://" + request.Host + request.Path
	req, err := http.NewRequest(request.Method, urlStr, body)
	if err != nil {
		return models.InternalError
	}

	for key, value := range request.Header {
		req.Header.Add(key, value)
	}

	inputSource, err := os.Open("./params")
	if err != nil {
		return models.NotFound
	}
	defer func(inputSource *os.File) {
		_ = inputSource.Close()
	}(inputSource)

	scanner := bufio.NewScanner(inputSource)
	var params []string
	for scanner.Scan() {
		params = append(params, scanner.Text())
	}

	randString := utils.RandStringRunes(rand.Intn(6))
	for _, param := range params {
		resReqWithParam := req
		query := resReqWithParam.URL.Query()
		query.Add(param, randString)
		resReqWithParam.URL.RawQuery = query.Encode()

		response, err := http.DefaultTransport.RoundTrip(resReqWithParam)
		if err != nil {
			continue
		}

		if response.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(response.Body)
			if err != nil {
				continue
			}

			bodyString := string(bodyBytes)
			if strings.Contains(bodyString, param) {
				fmt.Println(param)
			}
		}
		_ = response.Body.Close()
	}
	return models.Okey
}
