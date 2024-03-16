package ingestion

type TSOption func(opts *TSOptions)

type TSOptions struct {
	Token     string
	Url       string
	SaveFuncs []TrsSaveFunc
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

func TSWithSaveFunc(sf TrsSaveFunc) TSOption {
	return func(opts *TSOptions) {
		opts.SaveFuncs = append(opts.SaveFuncs, sf)
	}
}
