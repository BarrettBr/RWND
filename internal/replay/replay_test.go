package replay_test

import (
	"errors"
	"io"
	"testing"

	"github.com/BarrettBr/RWND/internal/model"
	"github.com/BarrettBr/RWND/internal/replay"
)

type fakeStore struct {
	recCh chan model.Record
	errCh chan error
}

func (s *fakeStore) Stream() (<-chan model.Record, <-chan error) {
	return s.recCh, s.errCh
}

func TestReplay_NewStore(t *testing.T) {
	if _, err := replay.New(nil); err == nil {
		t.Fatalf("Expected error for nil store")
	}
}

func TestReplay_Step_ReturnsError(t *testing.T) {
	s := &fakeStore{
		recCh: make(chan model.Record, 1),
		errCh: make(chan error, 1),
	}
	e, err := replay.New(s)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	want := errors.New("Fake Error")
	s.errCh <- want
	close(s.errCh)
	close(s.recCh)

	if rec, err := e.Step(); err == nil || err.Error() != want.Error() || rec != nil {
		t.Fatalf("Expected error %v and nil record, got rec=%v err=%v", want, rec, err)
	}
}

func TestReplay_Step_EOFWhenNoRecords(t *testing.T) {
	s := &fakeStore{
		recCh: make(chan model.Record, 1),
		errCh: make(chan error),
	}
	e, err := replay.New(s)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// Provide exactly one record then close the record channel
	s.recCh <- model.Record{}
	close(s.recCh)
	close(s.errCh)

	if rec, err := e.Step(); err != nil || rec == nil {
		t.Fatalf("expected record on first Step, got rec=%v err=%v", rec, err)
	}

	if rec, err := e.Step(); err != io.EOF || rec != nil {
		t.Fatalf("expected io.EOF on second Step, got rec=%v err=%v", rec, err)
	}
}
