package scaner

import "http-proxy/internal/models"

type UseCase interface {
	AllRequests() ([]models.Request, models.StatusCode)
	GetRequest(id int) (models.Request, models.StatusCode)
}

type Repository interface {
	AllRequests() ([]models.Request, models.StatusCode)
	GetRequest(id int) (models.Request, models.StatusCode)
}
