package scaner

import "http-proxy/internal/models"

type UseCase interface {
	GetRequests() ([]models.Request, models.StatusCode)
}

type Repository interface {
	GetRequests() ([]models.Request, models.StatusCode)
}
