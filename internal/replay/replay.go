package replay

import (
	"fmt"
	"io"
	"strings"

	"github.com/BarrettBr/RWND/internal/model"
)

type Store interface {
    Stream() (<-chan model.Record, <-chan error)
}

type Engine struct {
    store Store

    recCh <-chan model.Record
    errCh <-chan error
    done bool
}

func New(store Store) (*Engine, error){
    if store == nil {
        return nil, fmt.Errorf("Store not defined")
    }
    return  &Engine{store: store}, nil
}

func (e *Engine) Run() error {
    recCh, errCh := e.store.Stream()

    for recCh != nil || errCh != nil {
        select {
        case rec, ok := <-recCh:
            if !ok {
                recCh = nil
                continue
            }
            printRecord("REPLAY", rec)
        case err, ok := <-errCh:
            if !ok {
                errCh = nil
                continue
            }
            return err
        }
    }
    
    return nil
}

func printRecord(prefix string, rec model.Record) {
	fmt.Printf("%s id=%d %s %s -> %d\n",
		prefix,
		rec.ID,
		rec.Request.Method,
		rec.Request.URL,
		rec.Response.Status,
	)

	if len(rec.Request.Body) > 0 {
		fmt.Println("BODY:")
		fmt.Println("  " + strings.ReplaceAll(string(rec.Request.Body), "\n", "\n  "))
	}
}

func (e *Engine) Step() error {
    if e.done {
        return io.EOF
    }

    if e.recCh == nil && e.errCh == nil {
        e.recCh, e.errCh = e.store.Stream()
    }

    select{
    case err, ok := <-e.errCh:
        if ok && err != nil {
            return err
        }
        return e.Step()
    case rec, ok := <-e.recCh:
        if !ok {
            e.done = true
            return io.EOF
        }
        printRecord("STEP", rec)
    }
    return nil
}

func (e *Engine) StepLoop() error {
	if e.recCh == nil && e.errCh == nil {
		e.recCh, e.errCh = e.store.Stream()
	}

	for {
		fmt.Print("Press Enter for next, q to quit > ")
		var s string
		_, _ = fmt.Scanln(&s)
		if s == "q" {
			return nil
		}

		if err := e.Step(); err != nil {
			if err == io.EOF {
				fmt.Println("Done")
				return nil
			}
			return err
		}
	}
}

func (e *Engine) Reset() {
    // Function used for TUI in future if wanting to restep from beginning
    e.recCh = nil
    e.errCh = nil
    e.done = false
}
