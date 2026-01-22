package proxy

import (
	"net/http"
	"net/url"

	"github.com/BarrettBr/RWND/internal/model"
)


type Logger interface {
	Log(model.Record)
}

type Config struct {
	ListenAddr string
	Target     *url.URL
	Logger     Logger
}

type Proxy struct {
    cfg Config
    server *http.Server
}

// TODO: Finish New / Run
func New(cfg Config) *Proxy {
    return &Proxy{cfg:cfg}
}

func (p *Proxy) Run() error {
    return nil
}
