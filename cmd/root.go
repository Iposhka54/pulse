package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Iposhka54/pulse/internal/collector"
	"github.com/Iposhka54/pulse/internal/server"
	"github.com/Iposhka54/pulse/internal/ui"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

const (
	minInterval = time.Second
	maxInterval = 7 * 24 * time.Hour
)

func Execute(version string) error {
	root := cobra.Command{
		Use:   "pulse",
		Short: "Essential tool for your cluster",
		Long: ui.RenderBanner() + `

Pulse - Professional Server Monitoring Tool

Monitor your servers with real-time metrics and smart notification.
Keep a finger on the pulse of your servers and Kubernetes cluster. One command, instant health.

Need help?
Telegram: @machine_go
Vk: kafka_starter
`,
		DisableSuggestions: true,
		SilenceErrors:      true,
		SilenceUsage:       true,
		CompletionOptions:  cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	root.AddCommand(newStatusCmd())
	root.AddCommand(newServeCmd(version))
	return root.Execute()
}

func newStatusCmd() *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current system metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := collector.SnapshotNow()
			if err != nil {
				return err
			}
			if asJSON {
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				if err = encoder.Encode(s); err != nil {
					return err
				}
				return nil
			}

			ui.PrintStatus(s)
			return nil
		},
	}

	cmd.Flags().BoolVarP(
		&asJSON,
		"json", "j",
		false,
		"Output metrics as JSON",
	)

	return cmd
}

func newServeCmd(version string) *cobra.Command {
	var addr string
	var interval time.Duration

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run agent with HTTP API",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateInterval(interval); err != nil {
				return err
			}
			cfg := server.Config{
				Addr:            addr,
				ShutdownTimeout: 10 * time.Second,
			}
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()
			c := collector.New(interval)

			errGroup, ctx := errgroup.WithContext(ctx)
			errGroup.Go(func() error {
				c.Run(ctx)
				return nil
			})
			srv := server.NewHttpServer(cfg, c, version)
			errGroup.Go(func() error {
				err := server.Start(ctx, srv)
				cancel()
				return err
			})

			return errGroup.Wait()
		},
	}

	cmd.Flags().StringVarP(
		&addr,
		"addr", "a",
		":9100",
		"Address for http server ",
	)
	cmd.Flags().DurationVarP(
		&interval,
		"interval", "i",
		5*time.Second,
		"Interval for get system metrics",
	)

	return cmd
}

func validateInterval(d time.Duration) error {
	if d < minInterval {
		return fmt.Errorf("interval must be >= %s, got %s", minInterval, d)
	}
	if d > maxInterval {
		return fmt.Errorf("interval must be <= %s, got %s", maxInterval, d)
	}
	return nil
}
