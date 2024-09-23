package config

import "github.com/google/uuid"

type Config[C any] struct {
	Directory     string `default:"~/.config/et"`
	CommandConfig C
}

type User struct {
	Id        uuid.UUID
	Namespace string
}

type AppConfig struct {
	HomeDirectory string
	ActiveConfig  ActiveConfig
}

type ActiveConfig struct {
	Name      string `yaml:"-"`
	GridUrl   string `yaml:"grid_url"`
	UserToken string `yaml:"user_token"`
}

func NewActiveConfig(name string) ActiveConfig {
	return ActiveConfig{
		Name: name,
	}
}
