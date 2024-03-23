package management

import (
	"context"
	"log/slog"

	"github.com/dovydasdo/clerk/pkg/ingestion"
)

type RMOption func(opts *RMOptions)
type RMOptions struct {
	Processors []ingestion.Processor
	Ctx        context.Context
	Logger     *slog.Logger
}

func GetRMOptions(setters ...RMOption) *RMOptions {
	options := &RMOptions{}
	for _, s := range setters {
		s(options)
	}

	return options
}

func RMWithProcessor(p ingestion.Processor) RMOption {
	return func(opts *RMOptions) {
		opts.Processors = append(opts.Processors, p)
	}
}

func RMWithCtx(c context.Context) RMOption {
	return func(opts *RMOptions) {
		opts.Ctx = c
	}
}

func RMWithLogger(l *slog.Logger) RMOption {
	return func(opts *RMOptions) {
		opts.Logger = l
	}
}
