package cli

import (
	"time"

	"github.com/BarrettBr/RWND/internal/config"
	"github.com/BarrettBr/RWND/internal/datastore"
	"github.com/BarrettBr/RWND/internal/logpath"
	"github.com/BarrettBr/RWND/internal/replay"
)

func runReplay(args []string) error {
	cfg, err := config.FromReplayArgs(args, config.Load())
	if err != nil {
		PrintHelp()
		return err
	}

	logPath, err := logpath.ResolveReplayPath(cfg.LogPath)
	if err != nil {
		return err
	}

	store, err := datastore.NewFileStore(logPath, 500*time.Millisecond)
	if err != nil {
		return err
	}
	engine, err := replay.New(store)
	if err != nil {
		return err
	}

	return engine.StepLoop()
}
