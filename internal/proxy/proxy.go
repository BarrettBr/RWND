package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/BarrettBr/RWND/internal/model"
)

type Logger interface {
	Log(model.Record)
}

type Options struct {
	ListenAddr string
	Target     *url.URL
	Logger     Logger
}

type Proxy struct {
	srv *http.Server
}

func New(opts Options) (*Proxy, error) {
	// Check Options
	if opts.ListenAddr == "" {
		opts.ListenAddr = ":8080"
	}
	if opts.Target == nil {
		return nil, fmt.Errorf("Target is required")
	}
	if opts.Logger == nil {
		return nil, fmt.Errorf("Logger is required")
	}

	rp := httputil.NewSingleHostReverseProxy(opts.Target)

	// Capture / Log the repsonse inside the same record that the request came from
	rp.ModifyResponse = func(resp *http.Response) error {
		// Get the responses stored capture and type assert it
		// Use empty struct so overlapping package keys don't mess with this as well as the empty struct being efficient / easy to deal with
		// Got the idea to use an empty struct as the key from here:
		// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
		capAny := resp.Request.Context().Value(captureKey{})
		cap, ok := capAny.(*capture)
		if !ok || cap == nil {
			return nil
		}

		// Capture body, close stream and then recreate it since it was a stream it will be gone upon read
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		_ = resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		cap.rec.Response.Status = resp.StatusCode
		cap.rec.Response.Headers = resp.Header.Clone()
		cap.rec.Response.Body = bodyBytes
		cap.rec.Timestamp = time.Now().UTC()

		opts.Logger.Log(cap.rec)

		return nil
	}

	// Capture errors and log them
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		capAny := r.Context().Value(captureKey{})
		cap, ok := capAny.(*capture)
		if ok && cap != nil {
			cap.rec.Response.Status = http.StatusBadGateway
			cap.rec.Response.Body = []byte(err.Error())
			cap.rec.Timestamp = time.Now().UTC()
			opts.Logger.Log(cap.rec)
		}
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}

	mux := http.NewServeMux()

	// Handle request logging
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read and copy body like we did with responses above
		// however request bodies are optional so we guard clause it
		var reqBody []byte
		if r.Body != nil {
			reqBody, _ = io.ReadAll(r.Body)
			_ = r.Body.Close()

			// Restore the request body so the upstream STILL receives it
			r.Body = io.NopCloser(bytes.NewReader(reqBody))
		}

		// Create a record
		var rec model.Record
		rec.Request.Method = r.Method
		reqURL := r.URL
		if !reqURL.IsAbs() {
			reqURL = opts.Target.ResolveReference(reqURL)
		}
		rec.Request.URL = reqURL.String()
		rec.Request.Headers = r.Header.Clone()
		if r.Host != "" {
			rec.Request.Headers.Set("Host", r.Host)
		}
		rec.Request.Body = reqBody

		// Attach record to context of the request
		cap := &capture{rec: rec}
		ctx := context.WithValue(r.Context(), captureKey{}, cap)
		rp.ServeHTTP(w, r.WithContext(ctx))
	})

	server := &http.Server{
		Addr:              opts.ListenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Proxy{srv: server}, nil
}

func (p *Proxy) Run() error {
	if p.srv == nil {
		return fmt.Errorf("Proxy Run: Server is nil")
	}
	fmt.Printf("rwnd proxy listening on %s -> %s\n", p.srv.Addr, "(target)")

	err := p.srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (p *Proxy) Shutdown(ctx context.Context) error {
	if p.srv == nil {
		return nil
	}

	return p.srv.Shutdown(ctx)
}

// Internal types used for context capture

type captureKey struct{}

type capture struct {
	rec model.Record
}
