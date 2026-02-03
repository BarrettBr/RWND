// Command rwnd is the RWND CLI entry point.
package main

import (
	"log"
	"os"

	"github.com/BarrettBr/RWND/internal/cli"
)

// main runs the CLI and exits on error.
func main() {
	// Call to internal/cli/root for now but later we will support TUI & frontend display
	if err := cli.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
