package usecase

import "http-proxy/internal/pkg/proxy"

type UseCase struct {
	repo proxy.Repository
}

func NewRepoUsecase(repo proxy.Repository) proxy.UseCase {
	return &UseCase{repo: repo}
}
