package http

import (
	"context"
	"encoding/json"
	"github.com/danilkompaniets/movieapp-microservice/metadata/pkg/model"
	"github.com/danilkompaniets/movieapp-microservice/movie/internal/gateway"
	"github.com/danilkompaniets/movieapp-microservice/pkg/discovery"
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

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", "http://"+addrs[rand.Intn(len(addrs)-1)]+"/metadata", nil)

	if err != nil {
		return nil, err
	}

	log.Println("Request URL:", req.URL)

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()
	response, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	}
	if response.StatusCode != http.StatusOK {
		return nil, err
	}
	var v *model.Metadata

	err = json.NewDecoder(response.Body).Decode(v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
