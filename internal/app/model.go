package app

import "github.com/google/uuid"

type UpConfig struct {
	DryRun    bool
	Spec      []byte
	Filename  string
	Variables map[string]string
}

type UpResult struct {
	AggregateId uuid.UUID
	Namespace   string
}
