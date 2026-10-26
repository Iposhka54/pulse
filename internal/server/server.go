package server

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Config struct {
	Addr            string
	ShutdownTimeout time.Duration
}

type Server interface {
	Run(ctx context.Context) error
}

func Start(ctx context.Context, srv Server) error {
	if err := srv.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("run: %w", err)
	}
	return nil
}
