package delivery

import (
	"http-proxy/internal/pkg/proxy"
	"net/http"
)

type Handler struct {
	uc proxy.UseCase
}

func NewForumHandler(ForumUseCase proxy.UseCase) *Handler {
	return &Handler{uc: ForumUseCase}
}

func (h *Handler) AllRequest(w http.ResponseWriter, r *http.Request) {
}
