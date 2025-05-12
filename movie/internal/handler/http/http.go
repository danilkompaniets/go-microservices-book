package httpHandler

import (
	"encoding/json"
	"github.com/danilkompaniets/movieapp-microservice/movie/internal/controller/movie"
	"net/http"
)

type Handler struct {
	controller *movie.Controller
}

func New(controller *movie.Controller) *Handler {
	return &Handler{
		controller: controller,
	}
}

func (h *Handler) GetMovieDetails(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	ctx := r.Context()

	res, err := h.controller.Get(ctx, id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
