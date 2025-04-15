package bus

import (
	"errors"
	"reflect"
	"sync"
)

type BusOption func(*sync.Map) error

func WithHandler[In, Out any](handler Handler[In, Out]) BusOption {
	return func(m *sync.Map) error {
		var in In
		tIn := reflect.TypeOf(in)
		key, err := getRegistryKey(tIn)
		if err != nil {
			return err
		}

		if _, ok := m.Load(key); ok {
			return errors.New("system.core:not-found-handler")
		} else {
			m.Store(key, handler)
			return nil
		}
	}
}

type HandlerBus struct {
	registry *sync.Map
}

func NewHandlerBus(options ...BusOption) (*HandlerBus, error) {
	mp := &sync.Map{}
	for _, opt := range options {
		if err := opt(mp); err != nil {
			return nil, err
		}
	}

	return &HandlerBus{mp}, nil
}
