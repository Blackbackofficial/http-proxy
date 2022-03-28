package delivery

import (
	"http-proxy/internal/pkg/repeater"
	"net/http"
)

type Handler struct {
	uc repeater.UseCase
}

func NewRepeaterHandler(RepeaterUseCase repeater.UseCase) *Handler {
	return &Handler{uc: RepeaterUseCase}
}

func (h *Handler) AllRequest(w http.ResponseWriter, r *http.Request) {
	//
	//forumS := models.Response{Slug: slug}
	//forumS, status := h.uc.GetForum(forumS)
	//utils.Response(w, status, forumS)
}
