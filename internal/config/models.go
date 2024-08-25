package config

type Config[C any] struct {
	Directory     string `default:"~/.config/et"`
	CommandConfig C
}

type AppConfig struct {
	HomeDirectory string
	ActiveConfig  string
}
