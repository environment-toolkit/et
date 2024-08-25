package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type Server interface {
	Port() int
	Serve(l net.Listener) error
	Shutdown(ctx context.Context) error
}

type server struct {
	*http.Server

	port int
}

func (s *server) Port() int {
	return s.port
}

func NewServer(ctx context.Context, handler http.Handler) (Server, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}

	srv := &http.Server{
		Handler: handler,
	}

	go func() {
		if err := srv.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("error serving: %v\n", err)
		}
	}()

	return &server{
		Server: srv,
		port:   l.Addr().(*net.TCPAddr).Port,
	}, nil
}
