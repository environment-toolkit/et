package app

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/environment-toolkit/et/internal/config"
	"github.com/environment-toolkit/et/internal/server"
)

type Servers struct {
	GridUrl string
	Servers []server.Server
}

func NewServers(ctx context.Context, appConfig config.AppConfig) (*Servers, error) {
	activeConfig := appConfig.ActiveConfig
	if len(activeConfig.GridUrl) != 0 {
		return &Servers{
			GridUrl: activeConfig.GridUrl,
			Servers: nil,
		}, nil
	}

	pubsub := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})
	dbFilename := filepath.Join(appConfig.HomeDirectory, fmt.Sprintf("%s.db", activeConfig.Name))

	gridHandler, err := server.NewGrid(ctx, dbFilename, pubsub)
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

	return &Servers{
		GridUrl: fmt.Sprintf("http://localhost:%d", gridServer.Port()),
		Servers: []server.Server{gridServer, actuatorServer},
	}, nil
}
