package repo

import (
	"bytes"
	"github.com/jackc/pgx"
	"http-proxy/internal/models"
	"http-proxy/internal/pkg/proxy"
	"http-proxy/internal/pkg/utils"
	"io/ioutil"
	"net/http"
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

	err = rp.Conn.QueryRow(`INSERT INTO request (method, scheme, host, path, header, body)
			values ($1, $2, $3, $4, $5, $6) RETURNING id`,
		r.Method,
		r.URL.Scheme,
		r.URL.Host,
		r.URL.Path,
		reqHeaders,
		string(reqBody)).Scan(&reqId)

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

	err = rp.Conn.QueryRow(`INSERT INTO response (request_id, code, message, header, body)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		reqId,
		resp.StatusCode,
		resp.Status[4:],
		respHeaders,
		respBody).Scan(&respId)

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
