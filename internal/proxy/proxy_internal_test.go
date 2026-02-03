package proxy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/BarrettBr/RWND/internal/model"
)

type captureLogger struct {
	recCh chan model.Record
}

func (l *captureLogger) Log(rec model.Record) {
	l.recCh <- rec
}

func TestProxy_LogsRecord(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test", "ok")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("pong"))
	}))
	defer target.Close()

	targetURL, err := url.Parse(target.URL)
	if err != nil {
		t.Fatalf("parse target: %v", err)
	}

	logger := &captureLogger{recCh: make(chan model.Record, 1)}
	pxy, err := New(Options{
		ListenAddr: ":0",
		Target:     targetURL,
		Logger:     logger,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString("ping"))
	rr := httptest.NewRecorder()

	pxy.srv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	select {
	case rec := <-logger.recCh:
		if rec.Request.Method != http.MethodPost {
			t.Fatalf("expected method POST, got %s", rec.Request.Method)
		}
		if rec.Request.URL != target.URL+"/test" {
			t.Fatalf("expected url %s, got %s", target.URL+"/test", rec.Request.URL)
		}
		if rec.Response.Status != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Response.Status)
		}
		if string(rec.Response.Body) != "pong" {
			t.Fatalf("expected body pong, got %q", string(rec.Response.Body))
		}
		if rec.Response.Headers.Get("X-Test") != "ok" {
			t.Fatalf("expected header X-Test=ok")
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for log record")
	}
}
