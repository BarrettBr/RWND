package logger

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/BarrettBr/RWND/internal/model"
)


type Logger struct {
    store Store // Datastore
    ch chan model.Record // Channel to send records to get logged
    nextID atomic.Uint64 // Used so we can always call to next records id, atomic so if we use this we can increment it and not worry about duplicate ids in log
}

// Store interface so regardless of FileStore / SQLiteStore / etc it will still be supported
type Store interface {
    Append(model.Record) error
}

// ------------

// TODO: Add close method so on app close it logs final ones in buffer instead of dropping

func New(store Store) *Logger {
    l := &Logger{
        store: store,
        ch: make(chan model.Record, 1024), // Buffered channel so non-blocking if logs grow quickly
    }
    go l.worker()
    return l
}

func (l *Logger) Log(rec model.Record) {
    rec.ID = l.nextID.Add(1)
    rec.Timestamp = time.Now().UTC()

    // TODO: We just send it for now, maybe count drop or if buffer is full deal with that
    select {
    case l.ch <- rec:
    default:
    }
}

func (l *Logger) worker() {
    // Async worker func called in New
    // Iterate through loggers channel and store
    for rec := range l.ch {
        err := l.store.Append(rec)
        if err != nil {
            log.Printf("Logger append error: %s", err) // TODO: Eventually handle error instead of logging (Store / Callback)
        }
    }
}
