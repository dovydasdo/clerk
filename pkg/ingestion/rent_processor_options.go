package ingestion

import (
	"context"
	"log/slog"

	"github.com/go-co-op/gocron/v2"
)

type RPOption func(opt *RPOptions)
type RPOptions struct {
	Sources       []Ingestor
	Store         Saver
	Logger        *slog.Logger
	State         StateTracker
	StateF        []RentStateFunc
	StateInitFunc func(s StateTracker) StateTracker
	Ctx           context.Context
	Id            string
	DumpInterval  string
	Scheduler     gocron.Scheduler
}

func GetRPOptions(opts ...RPOption) *RPOptions {
	rpo := &RPOptions{}
	for _, opt := range opts {
		opt(rpo)
	}

	return rpo
}

func RPWithSource(s ...Ingestor) RPOption {
	return func(opt *RPOptions) {
		opt.Sources = append(opt.Sources, s...)
	}
}

func RPWithSaver(s Saver) RPOption {
	return func(opt *RPOptions) {
		opt.Store = s
	}
}

func RPWithLogger(l *slog.Logger) RPOption {
	return func(opt *RPOptions) {
		opt.Logger = l
	}
}

func RPWithState(s StateTracker) RPOption {
	return func(opt *RPOptions) {
		opt.State = s
	}
}

func RPWithStateF(sf RentStateFunc) RPOption {
	return func(opt *RPOptions) {
		if opt.StateF == nil {
			opt.StateF = make([]RentStateFunc, 0)
		}
		opt.StateF = append(opt.StateF, sf)
	}
}

func RPWithStateInitF(initf func(s StateTracker) StateTracker) RPOption {
	return func(opt *RPOptions) {
		opt.StateInitFunc = initf
	}
}

func RPWithCtx(ctx context.Context) RPOption {
	return func(opt *RPOptions) {
		opt.Ctx = ctx
	}
}

func RPWithId(id string) RPOption {
	return func(opt *RPOptions) {
		opt.Id = id
	}
}

func RPWithDumpInterval(interval string) RPOption {
	return func(opt *RPOptions) {
		opt.DumpInterval = interval
	}
}

func RPWithScheduler(s gocron.Scheduler) RPOption {
	return func(opt *RPOptions) {
		opt.Scheduler = s
	}
}
