package app

import (
	"context"
	"errors"

	"github.com/environment-toolkit/et/internal/config"
	"github.com/environment-toolkit/et/internal/grid"
	"github.com/google/uuid"
)

type Manager interface {
	Up(ctx context.Context, config UpConfig) (*UpResult, error)
	Close(ctx context.Context) error
}

type manager struct {
	servers     *Servers
	gridManager grid.Manager
}

func (m *manager) Up(ctx context.Context, config UpConfig) (*UpResult, error) {
	aggregateId := uuid.New()
	namespace := "default"

	spec := string(config.Spec)

	if err := m.gridManager.CommandNewSpec(ctx, aggregateId, namespace, spec, config.Variables); err != nil {
		return nil, err
	}

	return &UpResult{
		AggregateId: aggregateId,
		Namespace:   namespace,
	}, nil
}

func (m *manager) Close(ctx context.Context) error {
	if m.servers == nil {
		return nil
	}

	var errs []error

	for _, s := range m.servers.Servers {
		if err := s.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func NewManager(ctx context.Context, appConfig config.AppConfig) (Manager, error) {
	// if we don't have a grid url then start the grid.
	servers, err := NewServers(ctx, appConfig)
	if err != nil {
		return nil, err
	}

	gridManager, err := grid.NewManager(servers.GridUrl)
	if err != nil {
		return nil, err
	}

	return &manager{
		servers:     servers,
		gridManager: gridManager,
	}, nil
}
