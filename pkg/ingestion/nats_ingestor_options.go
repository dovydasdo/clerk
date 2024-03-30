package ingestion

import (
	"context"
	"log/slog"
)

type NIOption func(*NIOptions)

type NIOptions struct {
	ServerUrl  string
	Subject    string
	StreamName string
	Logger     *slog.Logger
	Ctx        context.Context
	Name       string
}

func GetNIOptions(opts ...NIOption) *NIOptions {
	ni := &NIOptions{}
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

func NIWithName(name string) NIOption {
	return func(n *NIOptions) {
		n.Name = name
	}
}
