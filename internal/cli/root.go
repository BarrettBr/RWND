package cli

import (
	"fmt"
)

// Basic command Help pretty printer
func PrintHelp() {
    // TODO: Double check this aligns with expected behavior
    fmt.Println(`rwnd - local HTTP traffic recorder & replay tool

Usage:
  rwnd proxy  [options]   Start reverse proxy and record traffic
  rwnd replay [options]   Replay recorded traffic
  rwnd help               Show this help

Examples:
  rwnd proxy --listen :8080 --target http://localhost:3000
  rwnd replay --step`)
}

func Run(args []string) error {
    if len(args) == 0 {
        PrintHelp()
        return fmt.Errorf("No command specified")
    }

    switch args[0]{
    case "proxy":
        return runProxy(args[1:])
    case "replay":
        return runReplay(args[1:])
    case "help", "-h", "--help":
        PrintHelp()
        return nil
    default:
        PrintHelp()
        return fmt.Errorf("Unknown command: %s", args[0])
    }
}
