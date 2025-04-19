package rest

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/amir-mln/amdp-task/services/objects/cmd/rest/routers"
	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/cmd_upload"
	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/qry_meta"
	"github.com/amir-mln/amdp-task/services/objects/internal/repo"
	"github.com/amir-mln/amdp-task/services/objects/internal/storage"
	"github.com/amir-mln/amdp-task/system/core/bus"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type Configs struct {
	SigCh          <-chan os.Signal
	ErrCh          chan<- error
	Logger         *zap.Logger
	DB             *sql.DB
	MinIO          *minio.Client
	MinIOBucket    string
	HTTPServerAddr string
	srv            *http.Server
}

var (
	once           sync.Once
	configs        *Configs
	termInProgress atomic.Bool
)

func SetConfigs(cfg *Configs) {
	once.Do(func() {
		configs = cfg
	})
}

func Terminate() {
	if !termInProgress.CompareAndSwap(false, true) || configs.srv == nil || configs.Logger == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	configs.Logger.Info("Shutting down HTTP server")
	err := configs.srv.Shutdown(ctx)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		configs.Logger.Info("Faced an error while shutting down", zap.Error(err))
	}
}

func Run() {
	if configs == nil {
		panic("Called [Run] of [api] with nill configs; Did you forget to call [SetConfig]?")
	}

	repo := repo.NewDbRepository(configs.Logger, configs.DB)
	fs := storage.NewObjectStorage(configs.Logger, configs.MinIO, configs.MinIOBucket)
	busOpts := []bus.BusOption{
		bus.WithHandler(cmd_upload.NewUploadCmdHandler(configs.Logger, repo, fs)),
		bus.WithHandler(qry_meta.NewMetaQryHandler(configs.Logger, repo)),
	}
	bus, err := bus.NewHandlerBus(busOpts...)
	if err != nil {
		configs.Logger.Error("Creating handler bus failed", zap.Error(err))
		configs.ErrCh <- err
	}

	router := routers.NewRootRouter(configs.Logger, bus)
	configs.srv = &http.Server{Addr: configs.HTTPServerAddr, Handler: router.Router()}
	localErrCh := make(chan error)
	go func() {
		configs.Logger.Info("Starting HTTP server", zap.String("ADDR", configs.srv.Addr))
		err := configs.srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			configs.Logger.Error("HTTP server crashed", zap.Error(err))
			localErrCh <- err
		}
	}()

	select {
	case <-configs.SigCh:
		Terminate()
	case configs.ErrCh <- <-localErrCh:
	}
}
