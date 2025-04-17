package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	"github.com/Netflix/go-env"
	"github.com/amir-mln/amdp-task/services/objects/cmd/api/routers"
	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/cmd_upload"
	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/qry_meta"
	"github.com/amir-mln/amdp-task/services/objects/internal/repo"
	"github.com/amir-mln/amdp-task/services/objects/internal/storage"
	"github.com/amir-mln/amdp-task/system/core/bus"
	"github.com/amir-mln/amdp-task/system/drivers/logging"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type config struct {
	LOG_ENVIRONMENT     logging.Environment         `env:"LOG_ENVIRONMENT"`
	LOG_FILE_ENC_TYPE   logging.EncoderType         `env:"LOG_FILE_ENC_TYPE"`
	LOG_FILE_LVLF       logging.LevelFilter         `env:"LOG_FILE_LVLF"`
	LOG_FILE_LVL        logging.ZapLevelUnmarshaler `env:"LOG_FILE_LVL"`
	LOG_FILE_PATH       string                      `env:"LOG_FILE_PATH"`
	LOG_KAFKA_ENC_TYPE  logging.EncoderType         `env:"LOG_KAFKA_ENC_TYPE"`
	LOG_KAFKA_LVLF      logging.LevelFilter         `env:"LOG_KAFKA_LVLF"`
	LOG_KAFKA_LVL       logging.ZapLevelUnmarshaler `env:"LOG_KAFKA_LVL"`
	LOG_KAFKA_TOPIC     string                      `env:"LOG_KAFKA_TOPIC"`
	HTTP_SERVER_ADDR    string                      `env:"HTTP_SERVER_ADDR"`
	POSTGRES_DSN        string                      `env:"POSTGRES_DSN"`
	KAFKA_BROKERS       []string                    `env:"KAFKA_BROKERS"`
	MINIO_ENDPOINT      string                      `env:"MINIO_ENDPOINT"`
	MINIO_ROOT_USER     string                      `env:"MINIO_ROOT_USER"`
	MINIO_ROOT_PASSWORD string                      `env:"MINIO_ROOT_PASSWORD"`
	MINIO_USE_SSL       bool                        `env:"MINIO_USE_SSL"`
	MINIO_BUCKET_NAME   string                      `env:"MINIO_BUCKET_NAME"`
}

func createZapLogger(config config) (*zap.Logger, error) {
	f, err := os.OpenFile(config.LOG_FILE_PATH, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	opts := []logging.LoggingOption{
		logging.WithEnvironment(config.LOG_ENVIRONMENT),
		logging.WithEncoder(config.LOG_FILE_ENC_TYPE),
		logging.WithLevelFilter(config.LOG_FILE_LVLF),
		logging.WithZapLevel(zapcore.Level(config.LOG_FILE_LVL)),
	}
	fzc, err := logging.NewZapCoreFile(f, opts...)
	if err != nil {
		return nil, err
	}

	lp, err := sarama.NewAsyncProducer(config.KAFKA_BROKERS, logging.NewKafkaLoggingConfig())
	if err != nil {
		return nil, err
	}
	opts2 := []logging.LoggingOption{
		logging.WithEnvironment(config.LOG_ENVIRONMENT),
		logging.WithEncoder(config.LOG_KAFKA_ENC_TYPE),
		logging.WithLevelFilter(config.LOG_KAFKA_LVLF),
		logging.WithZapLevel(zapcore.Level(config.LOG_KAFKA_LVL)),
	}
	lzc, err := logging.NewZapCoreWS(logging.NewKafkaWriteSyncer(lp, config.LOG_KAFKA_TOPIC), opts2...)
	if err != nil {
		return nil, err
	}

	return zap.New(zapcore.NewTee(fzc, lzc)), nil
}

func openDatabase(logger *zap.Logger, config config) (*sql.DB, error) {
	logger.Info("Opening the database driver", zap.String("dsn", config.POSTGRES_DSN))
	db, err := sql.Open("pgx", config.POSTGRES_DSN)
	if err != nil {
		logger.Error("Opening the database failed", zap.Error(err))
		return nil, err
	}
	if err := db.PingContext(context.Background()); err != nil {
		logger.Error("Pinging the database failed", zap.Error(err))
		return nil, err
	}

	return db, nil
}

func createMinIOClient(logger *zap.Logger, config config) (*minio.Client, error) {
	logger.Info("Creating minio client")
	m, err := minio.New(
		config.MINIO_ENDPOINT,
		&minio.Options{
			Creds:  credentials.NewStaticV4(config.MINIO_ROOT_USER, config.MINIO_ROOT_PASSWORD, ""),
			Secure: config.MINIO_USE_SSL,
		},
	)
	if err != nil {
		logger.Error("Creating minio client failed", zap.Error(err))
		return nil, err
	}

	return m, nil
}

var logger *zap.Logger
var srv *http.Server
var termInProgress atomic.Bool

func Terminate() {
	if !termInProgress.CompareAndSwap(false, true) || srv == nil || logger == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logger.Info("Shutting down HTTP server")
	err := srv.Shutdown(ctx)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		logger.Info("Faced an error while shutting down", zap.Error(err))
	}
}

func Run(sigCh <-chan os.Signal, errCh chan<- error) {
	var config config
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		errCh <- err
		return
	}

	logger, err = createZapLogger(config)
	if err != nil {
		errCh <- err
		return
	}

	db, err := openDatabase(logger, config)
	if err != nil {
		errCh <- err
		return
	}

	minio, err := createMinIOClient(logger, config)
	if err != nil {
		errCh <- err
		return
	}

	repo := repo.NewDbRepository(logger, db)
	fs := storage.NewObjectStorage(logger, minio, config.MINIO_BUCKET_NAME)
	busOpts := []bus.BusOption{
		bus.WithHandler(cmd_upload.NewUploadCmdHandler(logger, repo, fs)),
		bus.WithHandler(qry_meta.NewMetaQryHandler(logger, repo)),
	}
	bus, err := bus.NewHandlerBus(busOpts...)
	if err != nil {
		logger.Error("Creating handler bus failed", zap.Error(err))
		errCh <- err
	}

	router := routers.NewRootRouter(logger, bus)
	srv = &http.Server{Addr: config.HTTP_SERVER_ADDR, Handler: router.Router()}
	localErrCh := make(chan error)
	go func() {
		logger.Info("Starting HTTP server", zap.String("ADDR", srv.Addr))
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server crashed", zap.Error(err))
			localErrCh <- err
		}
	}()

	select {
	case <-sigCh:
		Terminate()
	case errCh <- <-localErrCh:
	}
}
