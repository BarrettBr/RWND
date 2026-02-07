// Command rwnd is the RWND CLI entry point.
package main

import (
	"log"
	"os"

	"github.com/BarrettBr/RWND/internal/cli"
	"github.com/BarrettBr/RWND/internal/tui"
)

// main runs the CLI and exits on error.
func main() {
	// No args -> TUI. Otherwise use CLI.
	if len(os.Args) == 1 {
		if err := tui.Run(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := cli.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
