package replay_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
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

func TestReplay_Replay_RequiresAbsoluteURL(t *testing.T) {
	s := &fakeStore{
		recCh: make(chan model.Record, 1),
		errCh: make(chan error, 1),
	}
	e, err := replay.New(s)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	rec := model.Record{}
	rec.Request.Method = "GET"
	rec.Request.URL = "/relative"

	if got, err := e.Replay(rec); err == nil || got != nil {
		t.Fatalf("expected error for relative URL, got rec=%v err=%v", got, err)
	}
}

func TestReplay_Replay_SendsRequest(t *testing.T) {
	s := &fakeStore{
		recCh: make(chan model.Record, 1),
		errCh: make(chan error, 1),
	}
	e, err := replay.New(s)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test", "ok")
		_, _ = w.Write([]byte("pong"))
	}))
	defer ts.Close()

	rec := model.Record{}
	rec.Request.Method = "GET"
	rec.Request.URL = ts.URL + "/ping"

	got, err := e.Replay(rec)
	if err != nil {
		t.Fatalf("Replay: %v", err)
	}
	if got.Response.Status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, got.Response.Status)
	}
	if string(got.Response.Body) != "pong" {
		t.Fatalf("expected body pong, got %q", string(got.Response.Body))
	}
	if got.Response.Headers.Get("X-Test") != "ok" {
		t.Fatalf("expected header X-Test=ok")
	}
}
