package http

import (
	"context"
	"encoding/json"
	"github.com/danilkompaniets/movieapp-microservice/movie/internal/gateway"
	"github.com/danilkompaniets/movieapp-microservice/pkg/discovery"
	"github.com/danilkompaniets/movieapp-microservice/rating/pkg/model"
	"log"
	"math/rand"
	"net/http"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordType model.RecordType, id model.RecordId) (model.RatingValue, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest(http.MethodGet, "http://"+addrs[rand.Intn(len(addrs)-1)]+"/rating", nil)
	if err != nil {
		return 0, err
	}
	log.Println("Request URL:", req.URL)

	req = req.WithContext(ctx)

	values := req.URL.Query()
	values.Add("type", string(recordType))
	values.Add("id", string(id))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return 0, err
	}

	var rating model.RatingValue
	err = json.NewDecoder(resp.Body).Decode(&rating)
	if err != nil {
		return 0, err
	}

	return rating, nil
}

func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordId, recordType model.RecordType, rating *model.Rating) error {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, "http://"+addrs[rand.Intn(len(addrs)-1)]+"/rating", nil)

	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("type", string(recordType))
	values.Add("id", string(recordID))
	values.Add("userId", string(rating.UserID))
	values.Add("ratingValue", string(rating.Value))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return gateway.ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
