// Package app holds shared application workflows used by CLI and TUI.
package app

import (
	"context"
	"time"

	"github.com/BarrettBr/RWND/internal/config"
	"github.com/BarrettBr/RWND/internal/datastore"
	"github.com/BarrettBr/RWND/internal/logger"
	"github.com/BarrettBr/RWND/internal/logpath"
	"github.com/BarrettBr/RWND/internal/proxy"
)

const proxyShutdownTimeout = 10 * time.Second

// RunProxy starts the proxy and blocks until it exits or the context is canceled.
func RunProxy(ctx context.Context, cfg config.AppConfig) error {
	logPath, err := logpath.ResolveRecordPath(cfg.LogPath, cfg.ListenAddr, cfg.TargetURL)
	if err != nil {
		return err
	}

	store, err := datastore.NewFileStore(logPath, 500*time.Millisecond)
	if err != nil {
		return err
	}
	logr := logger.New(store)

	pxy, err := proxy.New(proxy.Options{
		ListenAddr: cfg.ListenAddr,
		Target:     cfg.TargetURL,
		Logger:     logr,
	})
	if err != nil {
		logr.Close()
		_ = store.Close()
		return err
	}

	runErrCh := make(chan error, 1)
	go func() { runErrCh <- pxy.Run() }()

	select {
	case err := <-runErrCh:
		logr.Close()
		_ = store.Close()
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), proxyShutdownTimeout)
		defer cancel()

		_ = pxy.Shutdown(shutdownCtx)
		logr.Close()
		storeErr := store.Close()
		runErr := <-runErrCh

		if storeErr != nil {
			return storeErr
		}
		return runErr
	}
}
