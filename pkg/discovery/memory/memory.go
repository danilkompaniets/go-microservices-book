package memory

import (
	"context"
	"github.com/danilkompaniets/movieapp-microservice/pkg/discovery"
	"sync"
	"time"
)

type serviceName string
type instanceID string

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

type Registry struct {
	sync.RWMutex
	serviceAddresses map[serviceName]map[instanceID]*serviceInstance
}

func NewRegistry() *Registry {
	return &Registry{
		serviceAddresses: make(map[serviceName]map[instanceID]*serviceInstance),
	}
}

func (r *Registry) Register(ctx context.Context, id instanceID, serviceName serviceName, hostPort string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddresses[serviceName]; !ok {
		r.serviceAddresses[serviceName] = make(map[instanceID]*serviceInstance)
	}
	r.serviceAddresses[serviceName][id] = &serviceInstance{
		hostPort:   hostPort,
		lastActive: time.Now(),
	}

	return nil
}

func (r *Registry) Deregister(ctx context.Context, serviceName serviceName, id instanceID) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddresses[serviceName][id]; !ok {
		return discovery.ErrNotFound
	}

	delete(r.serviceAddresses[serviceName], id)

	return nil
}

func (r *Registry) ServiceAddresses(ctx context.Context, serviceName serviceName) ([]string, error) {
	serviceInstances := r.serviceAddresses[serviceName]
	if len(r.serviceAddresses[serviceName]) == 0 {
		return []string{}, discovery.ErrNotFound
	}
	result := make([]string, len(serviceInstances))

	for _, instance := range serviceInstances {
		if instance.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		result = append(result, instance.hostPort)
	}

	return result, nil

}

func (r *Registry) ReportHealthyState(ctx context.Context, serviceName serviceName, id instanceID) error {
	r.RLock()
	defer r.RUnlock()
	instance, ok := r.serviceAddresses[serviceName][id]
	if !ok {
		return discovery.ErrNotFound
	}
	instance.lastActive = time.Now()
	return nil
}
