package server

import (
	"context"
	"net/http"

	actuatorHandler "github.com/environment-toolkit/actuator/handler"
	gridHandler "github.com/environment-toolkit/grid/handler"
	"github.com/spf13/viper"

	"github.com/go-apis/eventsourcing/es"
	"github.com/go-apis/utils/xservice"
)

func NewGrid(ctx context.Context, dbFilename string, pubsub es.MemoryBusPubSub) (http.Handler, error) {
	v := viper.New()
	v.Set("service", "grid")
	v.Set("version", "1.0.0")
	v.Set("data.type", "sqlite")
	v.Set("data.sqlite.file", dbFilename)
	v.Set("stream.type", "mpub")
	v.Set("stream.memory.topic", "et")

	svc, err := xservice.NewService(ctx, v)
	if err != nil {
		return nil, err
	}

	return gridHandler.NewHandler(ctx, svc, pubsub)
}

func NewActuator(ctx context.Context, pubsub es.MemoryBusPubSub) (http.Handler, error) {
	v := viper.New()
	v.Set("service", "actuator")
	v.Set("version", "1.0.0")
	v.Set("data.type", "sqlite")
	v.Set("data.sqlite.memory", "true")
	v.Set("stream.type", "mpub")
	v.Set("stream.memory.topic", "et")

	svc, err := xservice.NewService(ctx, v)
	if err != nil {
		return nil, err
	}

	return actuatorHandler.NewHandler(ctx, svc, pubsub)
}
