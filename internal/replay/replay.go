package replay

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/BarrettBr/RWND/internal/model"
)

type Store interface {
	Stream() (<-chan model.Record, <-chan error)
}

type Engine struct {
	store  Store
	client *http.Client

	recCh <-chan model.Record
	errCh <-chan error
	done  bool
}

func New(store Store) (*Engine, error) {
	// New initializes a replay engine for a given store.
	if store == nil {
		return nil, fmt.Errorf("Store not defined")
	}
	engine := &Engine{
		store:  store,
		client: &http.Client{Timeout: 30 * time.Second},
	}
	return engine, nil
}

func printRequestPretty(rec model.Record) {
	// printRequestPretty prints a request view.
	fmt.Printf("Request #%d\n", rec.ID)
	fmt.Printf("%s %s\n", rec.Request.Method, rec.Request.URL)
	printHeaders(rec.Request.Headers)
	printBody(rec.Request.Body)
}

func printResponsePretty(title string, resp struct {
	Status  int
	Headers http.Header
	Body    []byte
}) {
	// printResponsePretty prints a response view.
	fmt.Printf("%s\n", title)
	fmt.Printf("Status: %d\n", resp.Status)
	printHeaders(resp.Headers)
	printBody(resp.Body)
}

func printHeaders(headers http.Header) {
	// printHeaders prints headers in a sorted order.
	if len(headers) == 0 {
		return
	}
	fmt.Println("Headers:")
	keys := make([]string, 0, len(headers))
	for k := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, v := range headers[k] {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}
}

func printBody(body []byte) {
	// printBody prints a body with indentation
	if len(body) == 0 {
		return
	}
	fmt.Println("Body:")
	fmt.Println("  " + strings.ReplaceAll(string(body), "\n", "\n  "))
}

func (e *Engine) Step() (*model.Record, error) {
	// Step returns the next record, or io.EOF when the stream ends.
	if e.done {
		return nil, io.EOF
	}

	if e.recCh == nil && e.errCh == nil {
		e.recCh, e.errCh = e.store.Stream()
	}

	for {
		if e.recCh == nil && e.errCh == nil {
			e.done = true
			return nil, io.EOF
		}

		select {
		case err, ok := <-e.errCh:
			if !ok {
				e.errCh = nil
				continue
			}
			if err != nil {
				return nil, err
			}
		case rec, ok := <-e.recCh:
			if !ok {
				e.recCh = nil
				if e.errCh == nil {
					e.done = true
					return nil, io.EOF
				}
				continue
			}
			return &rec, nil
		}
	}
}

func (e *Engine) StepLoop() error {
	// StepLoop runs the prompt for stepping and replaying.
	if e.recCh == nil && e.errCh == nil {
		e.recCh, e.errCh = e.store.Stream()
	}

	var current *model.Record
	for {
		fmt.Print("Press Enter for next, r to replay, q to quit > ")
		var s string
		_, _ = fmt.Scanln(&s)
		if s == "q" {
			return nil
		}

		if s == "r" {
			if current == nil {
				fmt.Println("No record to replay yet")
				continue
			}
			replayed, err := e.Replay(*current)
			if err != nil {
				fmt.Printf("Replay error: %v\n", err)
				continue
			}
			printResponsePretty("Old Response", current.Response)
			fmt.Println("---")
			printResponsePretty("New Response", replayed.Response)
			continue
		}

		rec, err := e.Step()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Done")
				return nil
			}
			return err
		}

		current = rec
		printRequestPretty(*rec)
	}
}

func (e *Engine) Reset() {
	// Reset clears stream state so stepping can restart.
	e.recCh = nil
	e.errCh = nil
	e.done = false
}

func (e *Engine) Replay(rec model.Record) (*model.Record, error) {
	// Replay re-sends a recorded request and returns the new response.
	reqURL, err := url.Parse(rec.Request.URL)
	if err != nil {
		return nil, fmt.Errorf("Replay invalid request URL: %w", err)
	}
	if !reqURL.IsAbs() {
		return nil, fmt.Errorf("Replay requires absolute request URL")
	}

	body := bytes.NewReader(rec.Request.Body)
	req, err := http.NewRequest(rec.Request.Method, reqURL.String(), body)
	if err != nil {
		return nil, err
	}

	if rec.Request.Headers != nil {
		req.Header = rec.Request.Headers.Clone()
	}
	req.Header.Del("Content-Length")
	req.Header.Del("Host")
	req.Header.Del("Accept-Encoding")
	req.ContentLength = int64(len(rec.Request.Body))

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	replayed := rec
	replayed.Response.Status = resp.StatusCode
	replayed.Response.Headers = resp.Header.Clone()
	replayed.Response.Body = respBody
	replayed.Timestamp = time.Now().UTC()

	return &replayed, nil
}
