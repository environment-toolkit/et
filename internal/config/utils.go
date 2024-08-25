package config

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func ResolvePath(path string) (string, error) {
	p := path

	if strings.HasPrefix(p, "~") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		p = dir + p[1:]
	}

	return filepath.Abs(p)
}

func GetActiveConfig(directory string, defaultName string) (string, error) {
	activeConfigPath := filepath.Join(directory, "active_config")
	if _, err := os.Stat(activeConfigPath); err != nil {
		if os.IsNotExist(err) {
			return defaultName, nil
		}
		return "", err
	}

	data, err := os.ReadFile(activeConfigPath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
