package memory

import (
	"context"
	"github.com/danilkompaniets/movieapp-microservice/metadata/internal/repository"
	"github.com/danilkompaniets/movieapp-microservice/metadata/pkg/model"
	"sync"
)

type Repository struct {
	sync sync.RWMutex
	data map[string]*model.Metadata
}

func New() *Repository {
	return &Repository{
		data: make(map[string]*model.Metadata),
	}
}

func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {
	r.sync.RLock()
	defer r.sync.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}

	return m, nil
}

func (r *Repository) Put(_ context.Context, m *model.Metadata) error {
	r.sync.Lock()
	defer r.sync.Unlock()
	r.data[m.ID] = m
	return nil
}
