package cli

import (
	"github.com/BarrettBr/RWND/internal/app"
	"github.com/BarrettBr/RWND/internal/config"
)

func runReplay(args []string) error {
	cfg, err := config.FromReplayArgs(args, config.Load())
	if err != nil {
		PrintHelp()
		return err
	}

	return app.RunReplay(cfg)
}
