package adapter

import (
	"go.uber.org/zap"
)

type options struct {
	log *zap.SugaredLogger
}

type Option interface {
	apply(opts *options)
}

func parseOptions(opts []Option) *options {
	o := newDefaultOptions()
	for _, opt := range opts {
		opt.apply(o)
	}

	return o
}

func newDefaultOptions() *options {
	return &options{log: zap.NewNop().Sugar()}
}

// WithLog
func WithLog(logger *zap.Logger) *withLogOptions {
	return &withLogOptions{Logger: logger}
}

type withLogOptions struct {
	Logger *zap.Logger
}

func (w *withLogOptions) apply(opts *options) {
	if w.Logger != nil {
		opts.log = w.Logger.Sugar()
	} else {
		opts.log = zap.NewNop().Sugar()
	}
}
