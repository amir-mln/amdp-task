package routers

import (
	"net/http"

	"github.com/amir-mln/amdp-task/services/objects/cmd/rest/routers/objects"
	"github.com/amir-mln/amdp-task/system/core/bus"
	"go.uber.org/zap"
)

type RootRouter struct {
	logger *zap.Logger
	bus    *bus.HandlerBus
}

func NewRootRouter(l *zap.Logger, hb *bus.HandlerBus) *RootRouter {
	return &RootRouter{
		logger: l,
		bus:    hb,
	}
}

func (router *RootRouter) Router() http.Handler {
	mux := &http.ServeMux{}
	objr := objects.NewObjectRouter(router.logger, router.bus)

	mux.Handle("/api/v1/objects/", http.StripPrefix("/api/v1", objr.Router()))

	return mux
}
