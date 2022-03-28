package delivery

import (
	"github.com/gorilla/mux"
	"http-proxy/internal/models"
	"http-proxy/internal/pkg/scaner"
	"http-proxy/internal/pkg/scaner/middleware"
	"net/http"
	"strconv"
)

type Handler struct {
	uc scaner.UseCase
}

func NewRepeaterHandler(RepeaterUseCase scaner.UseCase) *Handler {
	return &Handler{uc: RepeaterUseCase}
}

// AllRequests /requests
func (h *Handler) AllRequests(w http.ResponseWriter, r *http.Request) {
	requests, status := h.uc.AllRequests()
	middleware.Response(w, status, requests)
}

// GetRequest /requests/{id}
func (h *Handler) GetRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids, found := vars["id"]
	if !found {
		middleware.Response(w, models.NotFound, nil)
		return
	}
	id, err := strconv.Atoi(ids)
	if err != nil {
		middleware.Response(w, models.InternalError, nil)
		return
	}

	requests, status := h.uc.GetRequest(id)
	middleware.Response(w, status, requests)
}

// RepeatRequest /repeat/{id}
func (h *Handler) RepeatRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids, found := vars["id"]
	if !found {
		middleware.Response(w, models.NotFound, nil)
		return
	}
	id, err := strconv.Atoi(ids)
	if err != nil {
		middleware.Response(w, models.InternalError, nil)
		return
	}

	requests, status := h.uc.RepeatRequest(id)
	middleware.Response(w, status, requests)
}

// Scan /scan/{id}
func (h *Handler) Scan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids, found := vars["id"]
	if !found {
		middleware.Response(w, models.NotFound, nil)
		return
	}
	id, err := strconv.Atoi(ids)
	if err != nil {
		middleware.Response(w, models.InternalError, nil)
		return
	}

	status := h.uc.Scan(id)
	middleware.Response(w, status, nil)
}
