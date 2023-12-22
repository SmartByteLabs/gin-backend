package logger

import (
	"context"

	"github.com/princeparmar/gin-backend.git/pkg/constant"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger interface {
	Info(msg string)
	Infof(msg string, args ...interface{})
	Warn(msg string)
	Warnf(msg string, args ...interface{})
	Error(err error, msg string)
	Errorf(err error, msg string, args ...interface{})
	FatalIfError(err error, msg string)
	Fatal(err error, msg string)

	WithField(key string, value interface{}) Logger
}

type logger struct {
	log zerolog.Logger
}

func New() Logger {
	return &logger{
		log: log.With().Logger(),
	}
}

func LoggerFromContext(ctx context.Context) Logger {
	if ctx == nil {
		return New()
	}

	if l, ok := ctx.Value(constant.CtxKey_Logger).(Logger); ok {
		return l
	}

	return New()
}

func NewContextWithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, constant.CtxKey_Logger, l)
}

func (l *logger) Info(msg string) {
	l.log.Info().Msg(msg)
}

func (l *logger) Infof(msg string, args ...interface{}) {
	l.log.Info().Msgf(msg, args...)
}

func (l *logger) Warn(msg string) {
	l.log.Warn().Msg(msg)
}

func (l *logger) Warnf(msg string, args ...interface{}) {
	l.log.Warn().Msgf(msg, args...)
}

func (l *logger) Error(err error, msg string) {
	l.log.Error().Err(err).Msg(msg)
}

func (l *logger) Errorf(err error, msg string, args ...interface{}) {
	l.log.Error().Err(err).Msgf(msg, args...)
}

func (l *logger) Fatal(err error, msg string) {
	l.log.Fatal().Err(err).Msg(msg)
}

func (l *logger) FatalIfError(err error, msg string) {
	if err != nil {
		l.log.Fatal().Err(err).Msg(msg)
	}
}

func (l *logger) WithField(key string, value interface{}) Logger {
	return &logger{
		log: l.log.With().Interface(key, value).Logger(),
	}
}
