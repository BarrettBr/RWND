package config_test

import (
    "testing"

    "github.com/BarrettBr/RWND/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
    cfg := config.Load()
    if cfg.ListenAddr != ":8080" {
        t.Fatalf("expected ListenAddr=:8080, got %q", cfg.ListenAddr)
    }
    if cfg.LogPath != ".rwnd/logs" {
        t.Fatalf("expected LogPath=.rwnd/logs, got %q", cfg.LogPath)
    }
}

func TestFromProxyArgs_MissingTarget(t *testing.T) {
    _, err := config.FromProxyArgs([]string{}, config.Load())
    if err == nil {
        t.Fatalf("Expected error when --target missing")
    }
}

func TestFromProxyArgs_InvalidTargetURL(t *testing.T) {
    _, err := config.FromProxyArgs([]string{"--target", "http://[::1"}, config.Load())
    if err == nil {
        t.Fatalf("Expected error for invalid target url")
    }
}

func TestFromProxyArgs_AppliesOverrides(t *testing.T) {
    cfg, err := config.FromProxyArgs([]string{
        "--target", "http://localhost:3000",
        "--listen", ":9999",
        "--log", "tmp/log.jsonl",
    }, config.Load())
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

    if cfg.ListenAddr != ":9999" {
        t.Fatalf("Expected ListenAddr=:9999, got %q", cfg.ListenAddr)
    }
    if cfg.TargetURL == nil || cfg.TargetURL.String() != "http://localhost:3000" {
        t.Fatalf("Expected TargetURL=http://localhost:3000, got %+v", cfg.TargetURL)
    }
    if cfg.LogPath != "tmp/log.jsonl" {
        t.Fatalf("Expected LogPath=tmp/log.jsonl, got %q", cfg.LogPath)
    }
}

func TestFromReplayArgs_AppliesLogOverride(t *testing.T) {
    cfg, err := config.FromReplayArgs([]string{"--log", "x.jsonl"}, config.Load())
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if cfg.LogPath != "x.jsonl" {
        t.Fatalf("Expected LogPath=x.jsonl, got %q", cfg.LogPath)
    }
}
