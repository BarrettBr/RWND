package cli

import (
	"fmt"

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

    store, err := datastore.NewFileStore(cfg.LogPath)
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

    return pxy.Run()
}
