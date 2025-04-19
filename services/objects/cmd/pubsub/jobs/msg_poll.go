package jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync"

	"github.com/IBM/sarama"
	"github.com/amir-mln/amdp-task/system/core/messaging"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MessagePoll struct {
	msgRepo  messaging.Repository
	producer sarama.AsyncProducer
	logger   *zap.Logger
}

func NewMessagePoll(l *zap.Logger, r messaging.Repository, p sarama.AsyncProducer) *MessagePoll {
	return &MessagePoll{
		msgRepo:  r,
		producer: p,
		logger:   l,
	}
}

func (mp *MessagePoll) ProcessMessages(ctx context.Context, lim uint) {
	tx, err := mp.msgRepo.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		mp.logger.Error("Failed to begin transaction", zap.Error(err))
		return
	}
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	msgs, err := mp.msgRepo.GetMessagesTx(ctx, tx, lim)
	if err != nil {
		mp.logger.Error("Failed to poll messages", zap.Error(err))
		return
	}
	if len(msgs) == 0 {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for _, msg := range msgs {
			meta, _ := json.Marshal(msg)
			k := msg.ID.String()
			kafkaMsg := &sarama.ProducerMessage{
				Topic:    "objects.outgoing",
				Metadata: meta,
				Key:      sarama.ByteEncoder(k),
				Value:    sarama.ByteEncoder(msg.Body),
			}

			select {
			case mp.producer.Input() <- kafkaMsg:
				wg.Add(1)
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case _ = <-mp.producer.Errors():
				wg.Done()
			case msg := <-mp.producer.Successes():
				go func() {
					defer wg.Done()
					k, _ := msg.Key.Encode()
					id, _ := uuid.ParseBytes(k)
					select {
					case <-ctx.Done():
						return
					default:
						_ = mp.msgRepo.DeleteMessageByIDTx(ctx, tx, id)
					}

				}()
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Wait()
}
