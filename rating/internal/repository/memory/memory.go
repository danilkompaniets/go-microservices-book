package memory

import (
	"context"
	"github.com/danilkompaniets/movieapp-microservice/rating/internal/repository"
	"github.com/danilkompaniets/movieapp-microservice/rating/pkg/model"
	"sync"
)

type Repository struct {
	data map[model.RecordType]map[model.RecordId][]model.Rating
	sync sync.RWMutex
}

func New() *Repository {
	return &Repository{
		data: make(map[model.RecordType]map[model.RecordId][]model.Rating),
	}
}

func (repo *Repository) Get(_ context.Context, id model.RecordId, recordType model.RecordType) ([]model.Rating, error) {
	repo.sync.RLock()
	defer repo.sync.RUnlock()
	res, ok := repo.data[recordType][id]
	if !ok {
		return nil, repository.ErrNotFound
	}

	if len(res) == 0 {
		return nil, repository.ErrNotFound
	}

	return res, nil
}

func (repo *Repository) Put(_ context.Context, recordType model.RecordType, id model.RecordId, rating *model.Rating) error {
	repo.sync.Lock()
	defer repo.sync.Unlock()
	repo.data[recordType][id] = append(repo.data[recordType][id], *rating)
	return nil
}
