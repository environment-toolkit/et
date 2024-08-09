package resolver

import (
	"context"
	"fmt"
	"strings"

	"github.com/environment-toolkit/et/pkg/models"
)

type Session interface {
	NewValue(value string) (*models.Value, error)
	Resolve() error
}

type session struct {
	ctx                  context.Context
	environmentVariables map[string]string

	envs    []*models.Value
	secrets []*models.Value
}

func (m *session) NewValue(token string) (*models.Value, error) {
	// we only support "env" and "secret" values
	before, after, ok := strings.Cut(token, ":")
	if !ok {
		return nil, fmt.Errorf("invalid token: %s", token)
	}

	key := models.ValueType(before)
	v := &models.Value{Key: after, ValueType: key}

	switch key {
	case models.Env:
		m.envs = append(m.envs, v)
		return v, nil
	case models.Secret:
		m.secrets = append(m.secrets, v)
		return v, nil
	default:
		return nil, fmt.Errorf("invalid value type: %s", key)
	}
}

func (m *session) Resolve() error {

	// resolve the environment variables.
	for _, v := range m.envs {
		if value, ok := m.environmentVariables[v.Key]; ok {
			v.Value = &value
		}
	}

	for _, v := range m.secrets {
		// resolve the secret value.
		value := "ssssh"
		v.Value = &value
	}

	return nil
}

func NewSession(ctx context.Context, environmentVariables map[string]string) Session {
	return &session{
		ctx:                  ctx,
		environmentVariables: environmentVariables,
	}
}
