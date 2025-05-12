package httpHandler

import (
	"encoding/json"
	"github.com/danilkompaniets/movieapp-microservice/rating/internal/controller/rating"
	"github.com/danilkompaniets/movieapp-microservice/rating/pkg/model"
	"net/http"
	"strconv"
)

type Handler struct {
	controller *rating.Controller
}

func New(controller *rating.Controller) *Handler {
	return &Handler{
		controller: controller,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	id := model.RecordId(r.FormValue("id"))
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recordType := model.RecordType(r.FormValue("recordType"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		res, err := h.controller.GetAggregatedRating(r.Context(), recordType, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return

	case http.MethodPost:
		userID := model.UserID(r.FormValue("userId"))
		v, err := strconv.ParseFloat(r.FormValue("value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = h.controller.PutRecord(r.Context(), recordType, id, model.Rating{UserID: userID, Value: model.RatingValue(v)})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
