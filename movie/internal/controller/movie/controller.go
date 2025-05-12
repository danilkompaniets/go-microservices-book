package movie

import (
	"context"
	modelmetadata "github.com/danilkompaniets/movieapp-microservice/metadata/pkg/model"
	moviemodel "github.com/danilkompaniets/movieapp-microservice/movie/pkg/model"
	ratingmodel "github.com/danilkompaniets/movieapp-microservice/rating/pkg/model"
)

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordType ratingmodel.RecordType, id ratingmodel.RecordId) (ratingmodel.RatingValue, error)
	PutRating(ctx context.Context, recordID ratingmodel.RecordId, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*modelmetadata.Metadata, error)
}

type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{ratingGateway: ratingGateway, metadataGateway: metadataGateway}
}

func (c *Controller) Get(ctx context.Context, id string) (*moviemodel.MovieDetails, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	rating, err := c.ratingGateway.GetAggregatedRating(ctx, ratingmodel.RecordTypeMovie, ratingmodel.RecordId(id))

	if err != nil {
		return nil, err
	}

	movie := &moviemodel.MovieDetails{
		Rating:   float64(rating),
		Metadata: *metadata,
	}

	return movie, nil
}
