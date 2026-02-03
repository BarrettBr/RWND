// Package model defines shared structures for captured traffic.
package model

import (
	"net/http"
	"time"
)

// Record captures a request/response pair and metadata.
type Record struct {
	ID        uint64
	Timestamp time.Time

	Request struct {
		Method  string
		URL     string
		Headers http.Header
		Body    []byte
	}

	Response struct {
		Status  int
		Headers http.Header
		Body    []byte
	}
}
