package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/Netflix/go-env"
	"github.com/amir-mln/amdp-task/services/objects/cmd/pubsub"
	"github.com/amir-mln/amdp-task/services/objects/cmd/rest"
	"github.com/amir-mln/amdp-task/system/drivers/logging"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type envs struct {
	LogEnvironment      logging.Environment         `env:"LOG_ENVIRONMENT"`
	LogFileEncType      logging.EncoderType         `env:"LOG_FILE_ENC_TYPE"`
	LogFileLvlF         logging.LevelFilter         `env:"LOG_FILE_LVLF"`
	LofFileLvl          logging.ZapLevelUnmarshaler `env:"LOG_FILE_LVL"`
	LogFilePath         string                      `env:"LOG_FILE_PATH"`
	LogKafkaEncType     logging.EncoderType         `env:"LOG_KAFKA_ENC_TYPE"`
	LogKafkaLvlF        logging.LevelFilter         `env:"LOG_KAFKA_LVLF"`
	LogKafkaLvl         logging.ZapLevelUnmarshaler `env:"LOG_KAFKA_LVL"`
	LogKafkaTopic       string                      `env:"LOG_KAFKA_TOPIC"`
	HTTPServerAddr      string                      `env:"HTTP_SERVER_ADDR"`
	PostgresDSN         string                      `env:"POSTGRES_DSN"`
	KafkaBrokers        []string                    `env:"KAFKA_BROKERS"`
	MinIOEndpoint       string                      `env:"MINIO_ENDPOINT"`
	MinIORootUser       string                      `env:"MINIO_ROOT_USER"`
	MinIORootPassword   string                      `env:"MINIO_ROOT_PASSWORD"`
	MinIOUseSSL         bool                        `env:"MINIO_USE_SSL"`
	MinIOBucketName     string                      `env:"MINIO_BUCKET_NAME"`
	MessagePollInterval time.Duration               `env:"MESSAGE_POLL_INTERVAL"`
	MessagePollSize     uint                        `env:"MESSAGE_POLL_SIZE"`
	ShutdownTimeout     time.Duration               `env:"SHUTDOWN_TIMEOUT"`
}

func createZapLogger(e envs) (*zap.Logger, sarama.AsyncProducer, error) {
	f, err := os.OpenFile(e.LogFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	opts := []logging.LoggingOption{
		logging.WithEnvironment(e.LogEnvironment),
		logging.WithEncoder(e.LogFileEncType),
		logging.WithLevelFilter(e.LogFileLvlF),
		logging.WithZapLevel(zapcore.Level(e.LofFileLvl)),
	}
	fzc, err := logging.NewZapCoreFile(f, opts...)
	if err != nil {
		return nil, nil, err
	}

	lp, err := sarama.NewAsyncProducer(e.KafkaBrokers, logging.NewKafkaLoggingConfig())
	if err != nil {
		return nil, nil, err
	}
	opts2 := []logging.LoggingOption{
		logging.WithEnvironment(e.LogEnvironment),
		logging.WithEncoder(e.LogKafkaEncType),
		logging.WithLevelFilter(e.LogKafkaLvlF),
		logging.WithZapLevel(zapcore.Level(e.LogKafkaLvl)),
	}
	lzc, err := logging.NewZapCoreWS(logging.NewKafkaWriteSyncer(lp, e.LogKafkaTopic), opts2...)
	if err != nil {
		return nil, nil, err
	}

	logger := zap.New(zapcore.NewTee(fzc, lzc))
	return logger, lp, nil
}

func openDatabase(e envs) (*sql.DB, error) {
	db, err := sql.Open("pgx", e.PostgresDSN)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}

func createMinIOClient(e envs) (*minio.Client, error) {
	return minio.New(
		e.MinIOEndpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(e.MinIORootUser, e.MinIORootPassword, ""),
			Secure: e.MinIOUseSSL,
		},
	)
}

func createAsyncProd(e envs) (sarama.AsyncProducer, error) {
	host, _ := os.Hostname()
	if host == "" {
		host = "amdp-task.services.objects"
	}

	scfg := sarama.NewConfig()
	scfg.Version = sarama.V3_9_0_0
	scfg.Producer.Compression = sarama.CompressionSnappy
	scfg.Producer.RequiredAcks = sarama.WaitForAll
	scfg.Producer.Return.Successes = true
	scfg.Producer.Return.Errors = true
	return sarama.NewAsyncProducer(e.KafkaBrokers, scfg)
}

func broadcast[T any](count int, input <-chan T) []chan T {
	output := make([]chan T, count)
	for i := range count {
		output[i] = make(chan T)
	}

	go func() {
		defer func() {
			for _, outCh := range output {
				close(outCh)
			}
		}()

		for data := range input {
			for _, outCh := range output {
				outCh <- data
			}
		}
	}()

	return output
}

func main() {
	var envs envs
	_, err := env.UnmarshalFromEnviron(&envs)
	if err != nil {
		panic(err)
	}

	logger, logAp, err := createZapLogger(envs)
	if err != nil {
		panic(err)
	}

	db, err := openDatabase(envs)
	if err != nil {
		logger.Error("Failed to open database connection", zap.Error(err))
		os.Exit(1)
	}

	minio, err := createMinIOClient(envs)
	if err != nil {
		logger.Error("Failed to create MinIO client", zap.Error(err))
		os.Exit(1)
	}

	ap, err := createAsyncProd(envs)
	if err != nil {
		logger.Error("Failed to create async producer client", zap.Error(err))
		os.Exit(1)
	}

	close := func() {
		db.Close()
		ap.Close()
		logAp.Close()
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sigChs := broadcast(3, sigCh)
	apiErrCh, pubsubErrCh := make(chan error), make(chan error)
	rest.SetConfigs(&rest.Configs{
		SigCh:          sigChs[0],
		ErrCh:          apiErrCh,
		Logger:         logger,
		DB:             db,
		MinIO:          minio,
		MinIOBucket:    envs.MinIOBucketName,
		HTTPServerAddr: envs.HTTPServerAddr,
	})
	pubsub.SetConfigs(&pubsub.Configs{
		SigCh:               sigChs[1],
		ErrCh:               pubsubErrCh,
		Logger:              logger,
		DB:                  db,
		AsyncProd:           ap,
		MessagePollInterval: envs.MessagePollInterval,
		MessagePollSize:     envs.MessagePollSize,
	})

	go rest.Run()
	go pubsub.Run()
	select {
	case s := <-sigChs[len(sigChs)-1]:
		logger.Warn(
			"Signal Received",
			zap.Stringer("Signal", s),
			zap.Duration("Exit After", envs.ShutdownTimeout),
		)
		go close()
		time.AfterFunc(envs.ShutdownTimeout, func() { os.Exit(1) })
	case err := <-apiErrCh:
		logger.Error(
			"Error from api server",
			zap.Error(err),
			zap.Duration("Exit After", envs.ShutdownTimeout),
		)
		go pubsub.Terminate()
		go close()
		time.AfterFunc(envs.ShutdownTimeout, func() { os.Exit(1) })
	case err := <-pubsubErrCh:
		logger.Error(
			"Error from pubsub server",
			zap.Error(err),
			zap.Duration("Exit After", envs.ShutdownTimeout),
		)
		go rest.Terminate()
		go close()
		time.AfterFunc(envs.ShutdownTimeout, func() { os.Exit(1) })
	}
}
