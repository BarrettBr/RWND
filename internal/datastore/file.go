package datastore

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/BarrettBr/RWND/internal/model"
)

// TODO: Maybe extend with json encoder / writer so less calls down the line?
type FileStore struct {
    path string // Path of file
    mu sync.Mutex // Used for RW
    file *os.File // Used to hold the file itself
}

// ------------

func NewFileStore(path string) (*FileStore, error) {
    // Check if Directory exists and make it if not
    dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

    // Open file as append only and open / create if it doesn't exist
    f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return &FileStore{}, err
    }
    
    fs := &FileStore{path: path, file: f}
    return fs, nil
}

func (fs *FileStore) Append(rec model.Record) error {
    fs.mu.Lock()
    defer fs.mu.Unlock()

    data, err := json.Marshal(rec)
    if err != nil {
        return err
    }

    _, err = fs.file.Write(append(data, '\n'))
    return err
}

func (fs *FileStore) Stream() (<-chan model.Record, <-chan error) {
    // Iterate through logs streaming 1 at a time to the replay engine
    out := make(chan model.Record)
    errCh := make(chan error, 1)

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

func (fs *FileStore) Close() error {
    fs.mu.Lock()
    defer fs.mu.Unlock()

    if fs.file == nil {
        return nil
    }

    err := fs.file.Close()
    fs.file = nil
    return err
}
