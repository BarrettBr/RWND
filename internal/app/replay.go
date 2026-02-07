package app

import (
	"time"

	"github.com/BarrettBr/RWND/internal/config"
	"github.com/BarrettBr/RWND/internal/datastore"
	"github.com/BarrettBr/RWND/internal/logpath"
	"github.com/BarrettBr/RWND/internal/replay"
)

// RunReplay starts the replay engine and blocks until it exits.
func RunReplay(cfg config.AppConfig) error {
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
		_ = store.Close()
		return err
	}

	return engine.StepLoop()
}
