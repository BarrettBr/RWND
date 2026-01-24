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

    store := datastore.NewFileStore(cfg.LogPath)
    logr := logger.New(store)

    // TODO: Set up internal/proxy
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
