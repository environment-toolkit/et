package server

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func Test_It(t *testing.T) {
	ctx := context.Background()

	pubsub := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})

	h, err := NewGrid(ctx, "./test.db", pubsub)
	if err != nil {
		t.Fatal(err)
		return
	}

	srv, err := NewServer(ctx, h)
	if err != nil {
		t.Fatal(err)
		return
	}

	srv.Shutdown(ctx)
}
