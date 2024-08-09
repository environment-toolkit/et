package resolver

import (
	"os"
	"strings"
)

func EnvironmentVariables() map[string]string {
	m := make(map[string]string)
	environ := os.Environ()
	for _, env := range environ {
		key, value, ok := strings.Cut(env, "=")
		if !ok {
			continue
		}
		m[key] = value
	}
	return m
}
