package httpHandler

import (
	"encoding/json"
	"errors"
	"github.com/danilkompaniets/movieapp-microservice/metadata/internal/controller/metadata"
	"github.com/danilkompaniets/movieapp-microservice/metadata/internal/repository"
	"log"
	"net/http"
)

type Handler struct {
	ctrl *metadata.Controller
}

func New(ctrl metadata.Controller) *Handler {
	return &Handler{ctrl: &ctrl}
}

func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	m, err := h.ctrl.Get(ctx, id)

	if err != nil && errors.Is(err, repository.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to encode metadata: %v\n", err)
		return
	}
}
