package metadata

import (
	"context"
	"github.com/danilkompaniets/movieapp-microservice/metadata/internal/repository"
	"github.com/danilkompaniets/movieapp-microservice/metadata/pkg/model"
)

type metadataRepository interface {
	Get(context.Context, string) (*model.Metadata, error)
}
type Controller struct {
	Repo metadataRepository
}

func New(repo metadataRepository) *Controller {
	return &Controller{repo}
}
func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	m, err := c.Repo.Get(ctx, id)
	if err != nil {
		return nil, repository.ErrNotFound
	}

	return m, nil
}
