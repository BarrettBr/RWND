package proxy_test

import (
	"net/url"
	"sync"
	"testing"

	"github.com/BarrettBr/RWND/internal/model"
	"github.com/BarrettBr/RWND/internal/proxy"
)

type captureLogger struct {
	mu   sync.Mutex
	recs []model.Record
	ch   chan model.Record
}

func newCaptureLogger() *captureLogger {
	return &captureLogger{ch: make(chan model.Record, 10)}
}

func (l *captureLogger) Log(r model.Record) {
	l.mu.Lock()
	l.recs = append(l.recs, r)
	l.mu.Unlock()

	select {
	case l.ch <- r:
	default:
	}
}

func TestProxy_New_ValidateOptions(t *testing.T) {
	target, _ := url.Parse("http://example.com")
	logr := newCaptureLogger()

	if _, err := proxy.New(proxy.Options{ListenAddr: ":0", Target: nil, Logger: logr}); err == nil {
		t.Fatalf("Expected error when Target is nil")
	}

	if _, err := proxy.New(proxy.Options{ListenAddr: ":0", Target: target, Logger: nil}); err == nil {
		t.Fatalf("Expected error when Logger is nil")
	}

	if _, err := proxy.New(proxy.Options{ListenAddr: "", Target: target, Logger: logr}); err != nil {
		t.Fatalf("Expected no error when ListenAddr empty, got %v", err)
	}
}
