package cli

import (
	"time"

	"github.com/BarrettBr/RWND/internal/config"
	"github.com/BarrettBr/RWND/internal/datastore"
	"github.com/BarrettBr/RWND/internal/replay"
)


func runReplay(args []string) error {
    cfg, step, err := config.FromReplayArgs(args, config.Load())
    if err != nil {
        PrintHelp()
        return err
    }

    store, err := datastore.NewFileStore(cfg.LogPath, 500 * time.Millisecond)
    if err != nil {
        return err
    }
    engine, err := replay.New(store)
    if err != nil {
        return err
    }

    // Used for stepping forward one bit pausing and then letting step recall runReplay or what not
    if step {
        return engine.StepLoop()
    }

    return engine.Run()
}
