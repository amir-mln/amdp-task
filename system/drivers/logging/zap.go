package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getZapEncoder(o options) (zapcore.Encoder, error) {
	var cfg zapcore.EncoderConfig
	if o.environment == Development {
		cfg = zap.NewDevelopmentEncoderConfig()
	} else if o.environment == Production {
		cfg = zap.NewProductionEncoderConfig()
	} else {
		return nil, fmt.Errorf("") //TODO
	}

	if o.encoder == Console {
		return zapcore.NewConsoleEncoder(cfg), nil
	} else if o.encoder == JSON {
		return zapcore.NewJSONEncoder(cfg), nil
	} else {
		return nil, fmt.Errorf("") //TODO
	}
}

func getZapLevelEnabler(o options) (zapcore.LevelEnabler, error) {
	operators := map[LevelFilter]func(lvl1, lvl2 zapcore.Level) bool{
		Gte:   func(lvl1, lvl2 zapcore.Level) bool { return lvl1 >= lvl2 },
		Gt:    func(lvl1, lvl2 zapcore.Level) bool { return lvl1 > lvl2 },
		NotEq: func(lvl1, lvl2 zapcore.Level) bool { return lvl1 != lvl2 },
		Eq:    func(lvl1, lvl2 zapcore.Level) bool { return lvl1 == lvl2 },
		Lt:    func(lvl1, lvl2 zapcore.Level) bool { return lvl1 < lvl2 },
		Lte:   func(lvl1, lvl2 zapcore.Level) bool { return lvl1 <= lvl2 },
	}
	op, ok := operators[o.filter]
	if !ok {
		return nil, fmt.Errorf("") //TODO:
	}

	enabler := func(lvl zapcore.Level) bool {
		return op(lvl, o.zapLevel)
	}
	return zap.LevelEnablerFunc(enabler), nil
}

func NewZapCoreFile(file *os.File, opts ...LoggingOption) (zapcore.Core, error) {
	sync := zapcore.Lock(file)
	return NewZapCoreWS(sync, opts...)
}

func NewZapCoreWS(ws zapcore.WriteSyncer, opts ...LoggingOption) (zapcore.Core, error) {
	options := &options{}
	for _, f := range opts {
		err := f(options)
		if err != nil {
			return nil, err
		}
	}
	enabler, err := getZapLevelEnabler(*options)
	if err != nil {
		return nil, err
	}
	encoder, err := getZapEncoder(*options)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(encoder, ws, enabler)
	return core, nil
}

type kafkaWriteSyncer struct {
	producer sarama.AsyncProducer
	topic    string
}

func NewKafkaLoggingConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = false
	config.Producer.Return.Successes = false
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	config.Producer.Flush.Messages = 100
	config.Version = sarama.V3_9_0_0
	return config
}

func NewKafkaWriteSyncer(p sarama.AsyncProducer, topic string) zapcore.WriteSyncer {
	return &kafkaWriteSyncer{
		producer: p,
		topic:    topic,
	}
}

func (k *kafkaWriteSyncer) Write(p []byte) (n int, err error) {
	k.producer.Input() <- &sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.ByteEncoder(p),
	}
	return len(p), nil
}

// Sarama doesn't have a flush method natively, this is just a no-op
func (k *kafkaWriteSyncer) Sync() error {
	return nil
}
