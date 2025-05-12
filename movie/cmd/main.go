package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/danilkompaniets/movieapp-microservice/movie/internal/controller/movie"
	metadataHttp "github.com/danilkompaniets/movieapp-microservice/movie/internal/gateway/metadata/http"
	ratingHttp "github.com/danilkompaniets/movieapp-microservice/movie/internal/gateway/rating/http"
	httpHandler "github.com/danilkompaniets/movieapp-microservice/movie/internal/handler/http"
	"github.com/danilkompaniets/movieapp-microservice/pkg/discovery"
	"github.com/danilkompaniets/movieapp-microservice/pkg/discovery/consul"
	"log"
	"net/http"
	"time"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	instanceID := discovery.GenerateServiceId(serviceName)

	ctx := context.Background()
	err = registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			err := registry.ReportHealthyState(instanceID, serviceName)
			if err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	metadataGateway := metadataHttp.New(registry)
	ratingGateway := ratingHttp.New(registry)

	ctrl := movie.New(ratingGateway, metadataGateway)

	handler := httpHandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(handler.GetMovieDetails))
	url := fmt.Sprintf(":%d", port)
	err = http.ListenAndServe(url, nil)
	if err != nil {
		panic(err)
	}
}
