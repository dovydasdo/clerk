package ingestion

import "log/slog"

type TSOption func(opts *TSOptions)

type TSOptions struct {
	Token     string
	Url       string
	SaveFuncs []TrsSaveFunc
	Logger    *slog.Logger
}

func GetTSOptions(setters ...TSOption) TSOptions {
	opts := TSOptions{}

	for _, setter := range setters {
		setter(&opts)
	}

	return opts
}

func TSWithToken(t string) TSOption {
	return func(opts *TSOptions) {
		opts.Token = t
	}
}

func TSWithUrl(u string) TSOption {
	return func(opts *TSOptions) {
		opts.Url = u
	}
}

func TSWithLogger(l *slog.Logger) TSOption {
	return func(opts *TSOptions) {
		opts.Logger = l
	}
}

func TSWithSaveFunc(sf TrsSaveFunc) TSOption {
	return func(opts *TSOptions) {
		if opts.SaveFuncs == nil {
			opts.SaveFuncs = make([]TrsSaveFunc, 0)
		}
		opts.SaveFuncs = append(opts.SaveFuncs, sf)
	}
}
