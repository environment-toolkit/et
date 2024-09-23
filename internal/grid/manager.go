package grid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/motemen/go-loghttp"
)

type Manager interface {
	CommandNewSpec(ctx context.Context, aggregateId uuid.UUID, namespace string, data string, variables map[string]string) error
}

type manager struct {
	client ClientWithResponsesInterface
}

func (m *manager) CommandNewSpec(ctx context.Context, aggregateId uuid.UUID, namespace string, data string, variables map[string]string) error {
	body := CommandsNewSpec{
		AggregateId: aggregateId,
		Namespace:   namespace,
		Data:        data,
		Variables:   &variables,
	}
	resp, err := m.client.CommandNewSpecWithResponse(ctx, body)
	if err != nil {
		return err
	}
	if resp.JSON400 != nil {
		return fmt.Errorf("not found")
	}
	if resp.JSON401 != nil {
		return fmt.Errorf("not found")
	}

	return nil
}

func NewManager(url string, middlewares ...RequestEditorFn) (Manager, error) {
	cli := &http.Client{
		Transport: &loghttp.Transport{},
	}

	opts := make([]ClientOption, len(middlewares)+1)
	opts[0] = WithHTTPClient(cli)
	for i, m := range middlewares {
		opts[i+1] = WithRequestEditorFn(m)
	}

	client, err := NewClientWithResponses(url, opts...)
	if err != nil {
		return nil, err
	}

	return &manager{
		client: client,
	}, nil
}
