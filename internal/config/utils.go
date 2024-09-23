package config

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
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

func GetActiveConfigName(directory string, defaultName string) (string, error) {
	activeConfigPath := filepath.Join(directory, "active_config")
	if _, err := os.Stat(activeConfigPath); err != nil {
		if os.IsNotExist(err) {
			return defaultName, nil
		}
		return "", err
	}

	f, err := os.Open(activeConfigPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			continue
		}
		return text, nil
	}

	return defaultName, nil
}

func GetActiveConfig(directory string, defaultName string) (ActiveConfig, error) {
	activeConfigName, err := GetActiveConfigName(directory, defaultName)
	if err != nil {
		return NewActiveConfig(defaultName), err
	}

	activeConfig := NewActiveConfig(activeConfigName)

	activeConfigPath := filepath.Join(directory, fmt.Sprintf("%s.yaml", activeConfigName))
	if _, err := os.Stat(activeConfigPath); err != nil {
		if os.IsNotExist(err) {
			return activeConfig, nil
		}
		return activeConfig, err
	}

	data, err := os.ReadFile(activeConfigPath)
	if err != nil {
		return activeConfig, err
	}
	if err := yaml.Unmarshal(data, &activeConfig); err != nil {
		return activeConfig, err
	}
	return activeConfig, nil
}
