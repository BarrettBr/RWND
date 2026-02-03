package config

import (
	"flag"
	"fmt"
	"net/url"
)

type AppConfig struct {
	ListenAddr string // ":8080"
	TargetURL  *url.URL
	LogPath    string // ".rwnd/logs"
}

func Load() AppConfig {
	// Return a default Config struct and overwrite in arg call if specified overwrite
	return AppConfig{
		ListenAddr: ":8080",
		LogPath:    ".rwnd/logs",
	}
}

func FromProxyArgs(args []string, cfg AppConfig) (AppConfig, error) {
	// Function to parse arguments for the proxy command out
	// Logic in this function is referencing this Go by Example page
	// https://gobyexample.com/command-line-flags
	fs := flag.NewFlagSet("proxy", flag.ExitOnError)

	listen := fs.String(
		"listen",
		cfg.ListenAddr,
		"Address to listen on",
	)

	target := fs.String(
		"target",
		"",
		"Upstream target URL (required)",
	)

	logPath := fs.String(
		"log",
		cfg.LogPath,
		"Path to log file or directory",
	)

	if err := fs.Parse(args); err != nil {
		return AppConfig{}, err
	}

	if *target == "" {
		return AppConfig{}, fmt.Errorf("Missing required --target")
	}

	u, err := url.Parse(*target)
	if err != nil {
		return AppConfig{}, fmt.Errorf("Invalid target URL: %v", err)
	}

	cfg.ListenAddr = *listen
	cfg.TargetURL = u
	cfg.LogPath = *logPath

	return cfg, nil
}

func FromReplayArgs(args []string, cfg AppConfig) (AppConfig, error) {
	// Function to parse arguments for the replay command out
	fs := flag.NewFlagSet("replay", flag.ContinueOnError)
	fs.SetOutput(nil) // Set to nil so os.StdErr is used by default

	logPath := fs.String(
		"log",
		cfg.LogPath,
		"Path to log file or directory",
	)

	if err := fs.Parse(args); err != nil {
		return AppConfig{}, err
	}

	cfg.LogPath = *logPath
	return cfg, nil
}
