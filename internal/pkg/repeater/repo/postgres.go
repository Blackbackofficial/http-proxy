package repo

import (
	"github.com/jackc/pgx"
	"http-proxy/internal/pkg/repeater"
)

type repoPostgres struct {
	Conn *pgx.ConnPool
}

func NewRepoPostgres(Conn *pgx.ConnPool) repeater.Repository {
	return &repoPostgres{Conn: Conn}
}
