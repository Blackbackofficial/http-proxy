package repo

import (
	"bytes"
	"github.com/jackc/pgx"
	"http-proxy/internal/models"
	"http-proxy/internal/pkg/proxy"
	"http-proxy/internal/pkg/utils"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	InsertIntoResponse = "INSERT INTO response (request_id, code, message, cookies, header, body) values ($1, $2, $3, $4, $5, $6) RETURNING id"
	InsertIntoRequest  = "INSERT INTO request (method, scheme, host, path, cookies, header, body) values ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
)

type repoPostgres struct {
	Conn *pgx.ConnPool
}

func NewRepoPostgres(Conn *pgx.ConnPool) proxy.RepoProxy {
	return &repoPostgres{Conn: Conn}
}

func (rp *repoPostgres) SaveRequest(r *http.Request) (int, error) {
	var reqId int
	reqHeaders := utils.HeadToStr(r.Header)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	sCookies, err := utils.CookiesToString(r.Cookies())
	if err != nil {
		log.Printf("Error in parsing cookies:  %v", err)
	}

	row := rp.Conn.QueryRow(InsertIntoRequest, r.Method, r.URL.Scheme, r.URL.Host, r.URL.Path, sCookies, reqHeaders, string(reqBody))
	err = row.Scan(&reqId)
	if err != nil {
		return 0, err
	}

	return reqId, err
}

func (rp *repoPostgres) SaveResponse(reqId int, resp *http.Response) (models.Response, error) {
	var respId int
	respHeaders := utils.HeadToStr(resp.Header)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Response{}, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))

	sCookies, err := utils.CookiesToString(resp.Cookies())
	if err != nil {
		log.Printf("Error in parsing cookies:  %v", err)
	}

	row := rp.Conn.QueryRow(InsertIntoResponse, reqId, resp.StatusCode, resp.Status[4:], sCookies, respHeaders, respBody)
	err = row.Scan(&respId)
	if err != nil {
		return models.Response{}, err
	}

	response := models.Response{
		Id:        respId,
		RequestId: reqId,
		Code:      resp.StatusCode,
		Message:   resp.Status[4:],
		Header:    utils.StrToHeader(respHeaders),
		Body:      string(respBody),
	}

	return response, err
}
