package repo

import (
	"github.com/jackc/pgx"
	"http-proxy/internal/models"
	"http-proxy/internal/pkg/scaner"
)

const (
	SelectAllFromRequest = "SELECT id, method, scheme, host, path FROM request;"
)

type repoPostgres struct {
	Conn *pgx.ConnPool
}

func NewRepoPostgres(Conn *pgx.ConnPool) scaner.Repository {
	return &repoPostgres{Conn: Conn}
}

func (r *repoPostgres) GetRequests() ([]models.Request, models.StatusCode) {
	requests := make([]models.Request, 0)
	rows, err := r.Conn.Query(SelectAllFromRequest)
	if err != nil {
		return requests, models.NotFound
	}
	defer rows.Close()

	for rows.Next() {
		req := models.Request{}
		err = rows.Scan(&req.Id, &req.Method, &req.Scheme, &req.Host, &req.Path)
		if err != nil {
			return requests, models.InternalError
		}
		requests = append(requests, req)
	}
	return requests, models.Okey
}
