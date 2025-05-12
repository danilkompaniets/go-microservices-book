package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/danilkompaniets/movieapp-microservice/metadata/internal/controller/metadata"
	"github.com/danilkompaniets/movieapp-microservice/metadata/internal/handler/http"
	"github.com/danilkompaniets/movieapp-microservice/metadata/internal/repository/memory"
	"github.com/danilkompaniets/movieapp-microservice/pkg/discovery"
	"github.com/danilkompaniets/movieapp-microservice/pkg/discovery/consul"
	"log"
	"net/http"
	"time"
)

const serviceName = "metadata"

func main() {
	var port int

	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.Parse()

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	instanceID := discovery.GenerateServiceId(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.
				ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	log.Println("Starting metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
	handler := httpHandler.New(*ctrl)
	http.Handle("/metadata", http.HandlerFunc(handler.GetMetadata))
	url := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(url, nil); err != nil {
		panic(err)
	}
}
