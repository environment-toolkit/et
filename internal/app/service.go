package app

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/environment-toolkit/et/internal/config"
	"github.com/environment-toolkit/et/internal/grid"
	"github.com/environment-toolkit/et/internal/server"
	"github.com/go-apis/utils/xes"
	"github.com/google/uuid"
)

type UpResult struct {
	AggregateId uuid.UUID
	Namespace   string
}

type Manager interface {
	Up(ctx context.Context, spec string, variables map[string]string) (*UpResult, error)
	Close(ctx context.Context) error
}

type manager struct {
	servers     []server.Server
	gridManager grid.Manager
}

func (m *manager) Up(ctx context.Context, spec string, variables map[string]string) (*UpResult, error) {
	aggregateId := uuid.New()
	namespace := "default"

	if err := m.gridManager.CommandNewSpec(ctx, aggregateId, namespace, spec, variables); err != nil {
		return nil, err
	}

	return &UpResult{
		AggregateId: aggregateId,
		Namespace:   namespace,
	}, nil
}

func (m *manager) Close(ctx context.Context) error {
	var errs []error

	for _, s := range m.servers {
		if err := s.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func NewManager(ctx context.Context, appConfig config.AppConfig) (Manager, error) {
	// start the servers if needed.
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})

	signKey := RandStringRunes(24)
	dbFilename := filepath.Join(appConfig.HomeDirectory, fmt.Sprintf("%s.db", appConfig.ActiveConfig))

	gridHandler, err := server.NewGrid(ctx, signKey, dbFilename, pubsub)
	if err != nil {
		return nil, err
	}
	actuatorHandler, err := server.NewActuator(ctx, pubsub)
	if err != nil {
		return nil, err
	}

	gridServer, err := server.NewServer(ctx, gridHandler)
	if err != nil {
		return nil, err
	}
	actuatorServer, err := server.NewServer(ctx, actuatorHandler)
	if err != nil {
		return nil, err
	}

	servers := []server.Server{
		gridServer,
		actuatorServer,
	}

	security, err := xes.NewSecurity(signKey)
	if err != nil {
		return nil, err
	}

	gridUrl := fmt.Sprintf("http://localhost:%d", gridServer.Port())
	gridManager, err := grid.NewManager(gridUrl, security)
	if err != nil {
		return nil, err
	}

	return &manager{
		servers:     servers,
		gridManager: gridManager,
	}, nil
}
