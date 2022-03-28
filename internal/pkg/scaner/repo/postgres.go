package repo

import (
	"github.com/jackc/pgx"
	"http-proxy/internal/models"
	"http-proxy/internal/pkg/scaner"
	"http-proxy/internal/pkg/utils"
)

const (
	SelectAllFromRequest = "SELECT id, method, scheme, host, path, header, body FROM request;"
	SelectOneFromRequest = "SELECT id, method, scheme, host, path, header, body FROM request WHERE id = $1;"
)

type repoPostgres struct {
	Conn *pgx.ConnPool
}

func NewRepoPostgres(Conn *pgx.ConnPool) scaner.Repository {
	return &repoPostgres{Conn: Conn}
}

func (r *repoPostgres) AllRequests() ([]models.Request, models.StatusCode) {
	requests := make([]models.Request, 0)
	rows, err := r.Conn.Query(SelectAllFromRequest)
	if err != nil {
		return requests, models.NotFound
	}
	defer rows.Close()

	for rows.Next() {
		req := models.Request{}
		header := ""
		err = rows.Scan(&req.Id, &req.Method, &req.Scheme, &req.Host, &req.Path, &header, &req.Body)
		if err != nil {
			return requests, models.InternalError
		}
		req.Header = utils.StrToHeader(header)
		requests = append(requests, req)
	}
	return requests, models.Okey
}

func (r *repoPostgres) GetRequest(id int) (models.Request, models.StatusCode) {
	req := models.Request{}
	header := ""

	row := r.Conn.QueryRow(SelectOneFromRequest, id)
	err := row.Scan(&req.Id, &req.Method, &req.Scheme, &req.Host, &req.Path, &header, &req.Body)
	if err != nil {
		return models.Request{}, models.NotFound
	}

	req.Header = utils.StrToHeader(header)
	return req, models.Okey
}
