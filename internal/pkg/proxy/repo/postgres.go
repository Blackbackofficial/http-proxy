package repo

import (
	"github.com/jackc/pgx"
	"http-proxy/internal/pkg/proxy"
)

type repoPostgres struct {
	Conn *pgx.ConnPool
}

func NewRepoPostgres(Conn *pgx.ConnPool) proxy.Repository {
	return &repoPostgres{Conn: Conn}
}
