package datastore

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BarrettBr/RWND/internal/model"
)

type FileStore struct {
    path string // Path of file
    mu sync.Mutex // Used for RW
    file *os.File // Used to hold the file itself
    buf *bufio.Writer // Hold a buffered writer to lower total writes
    enc *json.Encoder // Used to encoder / feed to buffer

    flushInterval time.Duration
    stopFlush     chan struct{}
    stopOnce      sync.Once
}

// ------------

func NewFileStore(path string, flushInterval time.Duration) (*FileStore, error) {
    // Check if Directory exists and make it if not
    dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

    // Open file as append only and open / create if it doesn't exist
    f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    const memorySize = 64 * 1024

    buf := bufio.NewWriterSize(f, memorySize)
    enc := json.NewEncoder(buf)
    
    fs := &FileStore{
        path:path, 
        file:f,
        buf:buf,
        enc:enc,
        flushInterval: flushInterval,
        stopFlush:     make(chan struct{}),
    }

    fs.startFlushLoop()

    return fs, nil
}

func (fs *FileStore) Append(rec model.Record) error {
    fs.mu.Lock()
    defer fs.mu.Unlock()

    if fs.file == nil || fs.enc == nil {
        return os.ErrClosed
    }

    return fs.enc.Encode(rec)
}

func (fs *FileStore) Stream() (<-chan model.Record, <-chan error) {
    // Iterate through logs streaming 1 at a time to the replay engine
    out := make(chan model.Record)
    errCh := make(chan error, 1)

    // Push bufferred writes to the file before reading just to ensure it has something to read
    fs.mu.Lock()
    if fs.file == nil || fs.buf == nil {
        fs.mu.Unlock()
        errCh <- os.ErrClosed
        close(out)
        return out, errCh
    }
    flushErr := fs.buf.Flush()
    fs.mu.Unlock()
    
    if flushErr != nil {
        errCh <- flushErr
        close(out)
        return out, errCh
    }

    // Anonymous function that runs in a seperate goroutine
    // this will stream out logs 1 at a time to the replay engine and clean up the channels upon exiting
    go func ()  {
        defer close(out)
        defer close(errCh)

        f, err := os.Open(fs.path)
        if err != nil {
            errCh <- err
            return
        }
        defer f.Close()

        decoder := json.NewDecoder(f)
        for {
            var rec model.Record
            if err := decoder.Decode(&rec); err != nil {
                if err == io.EOF {
                    return
                }
                errCh <- err
                return
            }

            out <- rec
        }
    }()

    return out, errCh
}

func (fs *FileStore) startFlushLoop() {
    if fs.flushInterval <= 0 {
        return
    }

    ticker := time.NewTicker(fs.flushInterval)
    go func() {
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                fs.mu.Lock()
                // If closed, exit
                if fs.file == nil || fs.buf == nil {
                    fs.mu.Unlock()
                    return
                }
                _ = fs.buf.Flush()
                fs.mu.Unlock()
            case <-fs.stopFlush:
                return
            }
        }
    }()
}

func (fs *FileStore) Close() error {
    fs.stopOnce.Do(func() { close(fs.stopFlush) })
    
    fs.mu.Lock()
    if fs.file == nil {
        fs.mu.Unlock()
        return nil
    }
    buf := fs.buf
    file := fs.file
    fs.file, fs.buf, fs.enc = nil, nil, nil
    fs.mu.Unlock()

    flushErr := buf.Flush()
    closeErr := file.Close()

    if flushErr != nil {
        return flushErr
    }
    return closeErr
}
