package config

import (
	"context"
	"os"
	"path"

	"github.com/mcuadros/go-defaults"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type CommandLine[C any] struct {
	Log    *logrus.Logger
	App    AppConfig
	Config Config[C]
}

func New[C any](ctx context.Context, commandConfig C) (*CommandLine[C], error) {
	config := Config[C]{
		CommandConfig: commandConfig,
	}
	defaults.SetDefaults(&config)

	resolve, err := ResolvePath(config.Directory)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(resolve, 0755); err != nil {
		return nil, err
	}

	// activeConfig
	activeConfig, err := GetActiveConfig(resolve, "default")
	if err != nil {
		return nil, err
	}

	appConfig := AppConfig{
		HomeDirectory: resolve,
		ActiveConfig:  activeConfig,
	}

	log := logrus.New()
	log.SetOutput(&lumberjack.Logger{
		Filename: path.Join(resolve, "logs", "et.log"),
		MaxAge:   1,
		Compress: true,
	})

	return &CommandLine[C]{
		Log:    log,
		App:    appConfig,
		Config: config,
	}, nil
}
