package datastore

import (
	"os"
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

// TODO: Finish File functions
func NewFileStore(path string) *FileStore {

    return &FileStore{}
}

func (fs *FileStore) Append(rec model.Record) error {
    return nil
}

func (fs *FileStore) Stream() (<-chan model.Record, <-chan error) {

    return make(<-chan model.Record), make(<-chan error)
}

func (fs *FileStore) Close() error {
   return nil
}
