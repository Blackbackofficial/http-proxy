package delivery

import (
	"http-proxy/internal/pkg/scaner"
	"http-proxy/internal/pkg/scaner/middleware"
	"net/http"
)

type Handler struct {
	uc scaner.UseCase
}

func NewRepeaterHandler(RepeaterUseCase scaner.UseCase) *Handler {
	return &Handler{uc: RepeaterUseCase}
}

func (h *Handler) AllRequest(w http.ResponseWriter, r *http.Request) {
	requests, status := h.uc.GetRequests()
	middleware.Response(w, status, requests)
}
