package logging

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

type options struct {
	encoder     EncoderType
	environment Environment
	zapLevel    zapcore.Level
	filter      LevelFilter
}

type LoggingOption func(opts *options) error

func WithEncoder(et EncoderType) LoggingOption {
	return func(opts *options) error {
		if !et.Valid() {
			return fmt.Errorf("")
		}
		opts.encoder = et
		return nil
	}
}

func WithEnvironment(ev Environment) LoggingOption {
	return func(opts *options) error {
		if !ev.Valid() {
			return fmt.Errorf("")
		}
		opts.environment = ev
		return nil
	}
}

func WithZapLevel(lvl zapcore.Level) LoggingOption {
	return func(opts *options) error {
		if !(zapcore.DebugLevel <= lvl && lvl <= zapcore.FatalLevel) {
			return fmt.Errorf("")
		}
		opts.zapLevel = lvl
		return nil
	}
}

func WithLevelFilter(f LevelFilter) LoggingOption {
	return func(opts *options) error {
		if !f.Valid() {
			return fmt.Errorf("")
		}
		opts.filter = f
		return nil
	}
}
