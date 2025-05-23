package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error
	Deregister(ctx context.Context, instanceID string, serviceName string) error
	ServiceAddresses(ctx context.Context, serviceName string) ([]string, error)
	ReportHealthyState(instanceID string, serviceName string) error
}

var ErrNotFound = errors.New("no service addresses found")

func GenerateServiceId(serviceName string) string {
	return fmt.Sprintf("%s-%s", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
