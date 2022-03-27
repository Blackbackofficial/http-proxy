package usecase

import "http-proxy/internal/pkg/repeater"

type UseCase struct {
	repo repeater.Repository
}

func NewRepoUsecase(repo repeater.Repository) repeater.UseCase {
	return &UseCase{repo: repo}
}
