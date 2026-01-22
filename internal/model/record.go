package model

import (
	"net/http"
	"time"
)

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
