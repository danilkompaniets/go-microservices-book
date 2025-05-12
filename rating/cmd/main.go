package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/danilkompaniets/movieapp-microservice/pkg/discovery"
	"github.com/danilkompaniets/movieapp-microservice/pkg/discovery/consul"
	"github.com/danilkompaniets/movieapp-microservice/rating/internal/controller/rating"
	httpHandler "github.com/danilkompaniets/movieapp-microservice/rating/internal/handler/http"
	"github.com/danilkompaniets/movieapp-microservice/rating/internal/repository/memory"
	"log"
	"net/http"
	"time"
)

const serviceName = "ratings"

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

	log.Println("Starting metadata service")
	repo := memory.New()
	ctrl := rating.New(repo)
	handler := httpHandler.New(ctrl)
	http.Handle("/ratings", http.HandlerFunc(handler.Handle))
	url := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(url, nil); err != nil {
		panic(err)
	}
}
