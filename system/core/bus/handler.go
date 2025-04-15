package bus

import (
	"context"
	"errors"
	"reflect"
)

type Handler[In, Out any] interface {
	Handle(context.Context, In) (Out, error)
}

func Handle[Out, In any](ctx context.Context, bus *HandlerBus, in In) (Out, error) {
	key, err := getRegistryKey(reflect.TypeOf(in))
	if err != nil {
		var o Out
		return o, err
	}

	iface, ok := bus.registry.Load(key)
	if !ok {
		var o Out
		return o, errors.New("system.core:not-found-handler")
	}

	handler, ok := iface.(Handler[In, Out])
	if !ok {
		var o Out
		return o, errors.New("system.core:not-found-handler")
	}

	return handler.Handle(ctx, in)
}
