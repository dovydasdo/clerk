package ingestion

import (
	"context"
	"log/slog"
)

type RPOption func(opt *RPOptions)
type RPOptions struct {
	Sources       []Ingestor
	Store         Saver
	Logger        *slog.Logger
	State         any
	StateF        []RentStateFunc
	StateInitFunc func(s any)
	Ctx           context.Context
	Id            string
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

func RPWithState(s any) RPOption {
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

func RPWithStateInitF(initf func(s any)) RPOption {
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
