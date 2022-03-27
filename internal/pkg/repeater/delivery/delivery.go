package delivery

import (
	"http-proxy/internal/pkg/repeater"
	"net/http"
)

type Handler struct {
	uc repeater.UseCase
}

func NewForumHandler(ForumUseCase repeater.UseCase) *Handler {
	return &Handler{uc: ForumUseCase}
}

func (h *Handler) AllRequest(w http.ResponseWriter, r *http.Request) {
}
