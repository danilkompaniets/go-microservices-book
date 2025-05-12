package model

import "github.com/danilkompaniets/movieapp-microservice/metadata/pkg/model"

type MovieDetails struct {
	Rating   float64        `json:"rating"`
	Metadata model.Metadata `json:"metadata"`
}
