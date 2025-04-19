package pubsub

import (
	"context"
	"database/sql"
	"os"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/amir-mln/amdp-task/services/objects/cmd/pubsub/jobs"
	"github.com/amir-mln/amdp-task/system/core/messaging"
	"go.uber.org/zap"
)

type Configs struct {
	SigCh               <-chan os.Signal
	ErrCh               chan<- error
	Logger              *zap.Logger
	DB                  *sql.DB
	AsyncProd           sarama.AsyncProducer
	MessagePollInterval time.Duration
	MessagePollSize     uint
}

var (
	once   sync.Once
	config *Configs
)

func SetConfigs(cfg *Configs) {
	once.Do(func() {
		config = cfg
	})
}

// this is just a no-op, the [pubsub] package doesn't have
// any termination logic
func Terminate() {}

func Run() {
	if config == nil {
		panic("Called [Run] of [pubsub] with nill configs; Did you forget to call [SetConfig]?")
	}

	repo := messaging.NewRepository(config.Logger, config.DB)
	pollJob := jobs.NewMessagePoll(config.Logger, repo, config.AsyncProd)

	// scheduler, err := gocron.NewScheduler()
	// if err != nil {
	// 	config.Logger.Error("Failed to create a scheduler", zap.Error(err))
	// 	config.ErrCh <- err
	// 	return
	// }

	// _, err = scheduler.NewJob(
	// 	gocron.DurationJob(config.MessagePollInterval),
	// 	gocron.NewTask(pollJob.ProcessMessages, config.MessagePollSize),
	// )
	// if err != nil {
	// 	config.Logger.Error("Failed to create a job", zap.Error(err))
	// 	config.ErrCh <- err
	// 	return
	// }

	// scheduler.Start()
	// defer scheduler.Shutdown()
	config.Logger.Info("Poll Info", zap.Duration("Duration", config.MessagePollInterval))
	tic := time.NewTicker(config.MessagePollInterval)
scheduler:
	for {
		tic.Reset(config.MessagePollInterval)

		select {
		case <-tic.C:
			ctx, cancel := context.WithTimeout(context.Background(), config.MessagePollInterval)
			pollJob.ProcessMessages(ctx, config.MessagePollSize)
			cancel()
		case <-config.SigCh:
			Terminate()
			break scheduler
		}
	}
}
