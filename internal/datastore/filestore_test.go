package datastore_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BarrettBr/RWND/internal/datastore"
	"github.com/BarrettBr/RWND/internal/model"
)

func TestFileStore_CreatesDirectoryAndFile(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "logs", "rwnd.jsonl")

    fs, err := datastore.NewFileStore(path, time.Second)
    if err != nil {
        t.Fatalf("NewFileStore: %v", err)
    }
    defer func() { _ = fs.Close() }()

    if _, err := os.Stat(path); err != nil {
        t.Fatalf("Expected file to exist, stat error: %v", err)
    }
}

func TestFileStore_Append_SameCount(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "logs", "log.jsonl")

    fs, err := datastore.NewFileStore(path, time.Second)
    if err != nil {
        t.Fatalf("NewFileStore: %v", err)
    }
    defer func() { _ = fs.Close() }()

    const N = 10
    for i := 0; i < N; i++ {
        if err := fs.Append(model.Record{}); err != nil {
            t.Fatalf("Append %d: %v", i, err)
        }
    }


    out, errCh := fs.Stream()

    got := 0
    for range out {
        got ++
    }
    if err := <-errCh; err != nil {
        t.Fatalf("Stream error: %v", err)
    }

    if got != N {
        t.Fatalf("Expected %d records and got %d", N, got)
    }
}
