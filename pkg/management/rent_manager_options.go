package management

import "github.com/dovydasdo/clerk/pkg/ingestion"

type RMOption func(opts *RMOptions)
type RMOptions struct {
	Processors []ingestion.Processor
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
