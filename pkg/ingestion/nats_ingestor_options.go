package ingestion

import (
	"context"
	"log/slog"
)

type NIOption func(*NIOptions)

type NIOptions struct {
	ServerUrl  string
	PFunc      Procedure
	DFunc      DumpFunc
	Subject    string
	StreamName string
	Logger     *slog.Logger
	Ctx        context.Context
	State      any
	Name       string
}

func GetNIOptions(opts ...NIOption) *NIOptions {
	ni := &NIOptions{
		Ctx:   context.Background(),
		State: make(map[string]int),
	}
	for _, opt := range opts {
		opt(ni)
	}
	return ni
}

func NIWithServerUrl(url string) NIOption {
	return func(n *NIOptions) {
		n.ServerUrl = url
	}
}

func NIWithPFunc(p Procedure) NIOption {
	return func(n *NIOptions) {
		n.PFunc = p
	}
}

func NIWithDFunc(d DumpFunc) NIOption {
	return func(n *NIOptions) {
		n.DFunc = d
	}
}

func NIWithSubject(sbj string) NIOption {
	return func(n *NIOptions) {
		n.Subject = sbj
	}
}

func NIWithStreamName(sname string) NIOption {
	return func(n *NIOptions) {
		n.StreamName = sname
	}
}

func NIWithLogger(l *slog.Logger) NIOption {
	return func(n *NIOptions) {
		n.Logger = l
	}
}

func NIWithCtx(ctx context.Context) NIOption {
	return func(n *NIOptions) {
		n.Ctx = ctx
	}
}

func NIWithState(state any) NIOption {
	return func(n *NIOptions) {
		n.State = state
	}
}

func NIWithName(name string) NIOption {
	return func(n *NIOptions) {
		n.Name = name
	}
}
