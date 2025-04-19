package logging

import (
	"fmt"

	"github.com/Netflix/go-env"
	"go.uber.org/zap/zapcore"
)

//go:generate stringer -type=EncoderType -output=cfg_enc_string.go
type EncoderType uint

const (
	Console EncoderType = iota
	JSON
)

var (
	dummyEnc EncoderType     = Console
	_        env.Unmarshaler = &dummyEnc
)

// [UnmarshalEnvironmentValue] implements [env.Unmarshaler]
func (enc *EncoderType) UnmarshalEnvironmentValue(data string) error {
	var err error
	switch data {
	case Console.String(), "":
		*enc = Console
	case JSON.String():
		*enc = JSON
	default:
		err = fmt.Errorf("") // TODO:Customer Error
	}
	return err
}

func (enc EncoderType) Valid() bool {
	return Console <= enc && enc <= JSON

}

//go:generate stringer -type=Environment -output=cfg_env_string.go
type Environment uint

const (
	Development Environment = iota
	Production
)

var (
	dummyEnv Environment     = Development
	_        env.Unmarshaler = &dummyEnv
)

func (e *Environment) UnmarshalEnvironmentValue(data string) error {
	var err error
	switch data {
	case Development.String(), "":
		*e = Development
	case Production.String():
		*e = Production
	default:
		err = fmt.Errorf("") // TODO:Customer Error
	}
	return err
}

func (e Environment) Valid() bool {
	return Development <= e && e <= Production
}

//go:generate stringer -type=LevelFilter -output=cfg_lvlf_string.go
type LevelFilter uint

const (
	Gte LevelFilter = iota
	Gt
	NotEq
	Eq
	Lt
	Lte
)

var (
	dummyLvlF LevelFilter     = Gte
	_         env.Unmarshaler = &dummyLvlF
)

func (f *LevelFilter) UnmarshalEnvironmentValue(data string) error {
	var err error
	switch data {
	case Lt.String():
		*f = Lt
	case Lte.String():
		*f = Lte
	case NotEq.String():
		*f = NotEq
	case Eq.String():
		*f = Eq
	case Gt.String():
		*f = Gt
	case Gte.String(), "":
		*f = Gte
	default:
		err = fmt.Errorf("") // TODO:Customer Error
	}
	return err
}

func (f LevelFilter) Valid() bool {
	return Gte <= f && f <= Lte
}

type ZapLevelUnmarshaler zapcore.Level

var (
	dummyZLU ZapLevelUnmarshaler = ZapLevelUnmarshaler(zapcore.DebugLevel)
	_        env.Unmarshaler     = &dummyZLU
)

func (z *ZapLevelUnmarshaler) UnmarshalEnvironmentValue(data string) error {
	if data == "" {
		*z = ZapLevelUnmarshaler(zapcore.DebugLevel)
		return nil
	}

	l, err := zapcore.ParseLevel(data)
	*z = ZapLevelUnmarshaler(l)
	return err
}
