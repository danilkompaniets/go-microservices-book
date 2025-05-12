package rating

import (
	"context"
	"github.com/danilkompaniets/movieapp-microservice/rating/internal/repository"
	"github.com/danilkompaniets/movieapp-microservice/rating/pkg/model"
)

type RatingsRepo interface {
	Get(_ context.Context, id model.RecordId, recordType model.RecordType) ([]model.Rating, error)
	Put(_ context.Context, recordType model.RecordType, id model.RecordId, rating *model.Rating) error
}

type Controller struct {
	repo RatingsRepo
}

func New(repo RatingsRepo) *Controller {
	return &Controller{
		repo: repo,
	}
}

func (c *Controller) PutRecord(ctx context.Context, recordType model.RecordType, recordId model.RecordId, r model.Rating) error {
	err := c.repo.Put(ctx, recordType, recordId, &r)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) GetAggregatedRating(ctx context.Context, recordType model.RecordType, id model.RecordId) (float64, error) {
	m, err := c.repo.Get(ctx, id, recordType)
	if err != nil {
		return 0, repository.ErrNotFound
	}

	sum := float64(0)
	for _, rating := range m {
		sum += float64(rating.Value)
	}

	return sum / float64(len(m)), nil
}
