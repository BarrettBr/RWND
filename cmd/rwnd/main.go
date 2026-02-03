package main

import (
	"log"
	"os"

	"github.com/BarrettBr/RWND/internal/cli"
)

func main() {
	// Call to internal/cli/root for now but later we will support TUI & frontend display
	if err := cli.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
