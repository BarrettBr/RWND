package logger_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/BarrettBr/RWND/internal/logger"
	"github.com/BarrettBr/RWND/internal/model"
)

type fakeStore struct {
	count atomic.Int64
	wg    sync.WaitGroup
}

func (s *fakeStore) Append(model.Record) error {
	s.count.Add(1)
	s.wg.Done()
	return nil
}

func TestLogger_Close_DrainsLogs(t *testing.T) {
	var s fakeStore
	l := logger.New(&s)

	const N = 200
	s.wg.Add(N)

	for range N {
		l.Log(model.Record{})
	}

	// Setup a channel to signal when the logger drains all logs
	done := make(chan struct{})
	go func() {
		l.Close()
		close(done)
	}()

	// If it doesn't drain we fail the test after 2 seconds to give test an upper bound
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("logger.Close() timed out; likely not draining")
	}

	// Wait till all appended, just done differently from the close check above since it could later
	// return early and not actually add them all
	waitDone := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(waitDone)
	}()

	select {
	case <-waitDone:
	case <-time.After(2 * time.Second):
		t.Fatalf("store did not receive all appends")
	}

	// Check if the added number of items matches the expected number
	if got := s.count.Load(); got != N {
		t.Fatalf("expected %d appends, got %d", N, got)
	}
}
