package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BarrettBr/RWND/internal/config"
	"github.com/BarrettBr/RWND/internal/datastore"
	"github.com/BarrettBr/RWND/internal/logger"
	"github.com/BarrettBr/RWND/internal/proxy"
)


func runProxy(args []string) error {
    cfg, err := config.FromProxyArgs(args, config.Load())
    if err != nil {
        PrintHelp()
        return fmt.Errorf("Proxy args Error: %v", err)
    }

    store, err := datastore.NewFileStore(cfg.LogPath, 500 * time.Millisecond)
    if err != nil {
        return err
    }
    logr := logger.New(store)

    pxy, err := proxy.New(proxy.Options{
        ListenAddr: cfg.ListenAddr,
        Target: cfg.TargetURL,
        Logger: logr,
    })
    if err != nil {
        return err
    }

    // Capture signal and do graceful shutdown off of it
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
    defer signal.Stop(sigCh)

    // Start server and feed back into a channel it's error
    runErrCh := make(chan error, 1)
    go func(){runErrCh <- pxy.Run()}()

    // If exit early with error return it otherwise if provided a terminate call the shutdown / close functions
    select{
    case err := <-runErrCh:
        logr.Close()
        _ = store.Close()
        return err
    case sig := <-sigCh:
        fmt.Printf("\nReceived %s, shutting down...\n", sig)

        ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
        defer cancel()

        _ = pxy.Shutdown(ctx)
        logr.Close()
        storeErr := store.Close()
        runErr := <- runErrCh

        if storeErr != nil {
            return storeErr
        }
        return runErr
    }
}
