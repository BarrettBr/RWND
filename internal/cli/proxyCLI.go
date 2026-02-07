package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/BarrettBr/RWND/internal/app"
	"github.com/BarrettBr/RWND/internal/config"
)

func runProxy(args []string) error {
	cfg, err := config.FromProxyArgs(args, config.Load())
	if err != nil {
		PrintHelp()
		return fmt.Errorf("Proxy args Error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capture signal and do graceful shutdown off of it.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	runErrCh := make(chan error, 1)
	go func() { runErrCh <- app.RunProxy(ctx, cfg) }()

	select {
	case err := <-runErrCh:
		return err
	case sig := <-sigCh:
		fmt.Printf("\nReceived %s, shutting down...\n", sig)
		cancel()
		runErr := <-runErrCh
		if errors.Is(runErr, context.Canceled) {
			return nil
		}
		return runErr
	}
}
