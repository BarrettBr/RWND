package replay

import (
	"fmt"

	"github.com/BarrettBr/RWND/internal/model"
)

type Store interface {
    Append(model.Record) error
}

type Engine struct {
    store Store
}

func New(store Store) (*Engine, error){
    if store == nil {
        return &Engine{}, fmt.Errorf("Store not defined")
    }

    return  &Engine{store: store}, nil
}

func (e *Engine) Run() error {
    return nil
}

func (e *Engine) Step() error {

    return nil
}
