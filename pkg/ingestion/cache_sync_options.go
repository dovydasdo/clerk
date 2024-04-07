package ingestion

import "context"

type NKVCacheOption func(opts *NKVCacheOptions)
type NKVCacheOptions struct {
	Name    string
	Url     string
	Ctx     context.Context
	Subject string
}

func GetNKVCacheOptions(options ...NKVCacheOption) *NKVCacheOptions {
	opts := &NKVCacheOptions{}

	for _, opt := range options {
		opt(opts)
	}

	return opts
}

func NKVCacheWithName(n string) NKVCacheOption {
	return func(opts *NKVCacheOptions) {
		opts.Name = n
	}
}

func NKVCacheWithUrl(u string) NKVCacheOption {
	return func(opts *NKVCacheOptions) {
		opts.Url = u
	}
}

func NKVCacheWithCtx(c context.Context) NKVCacheOption {
	return func(opts *NKVCacheOptions) {
		opts.Ctx = c
	}
}

func NKVCacheWithSubject(s string) NKVCacheOption {
	return func(opts *NKVCacheOptions) {
		opts.Subject = s
	}
}
