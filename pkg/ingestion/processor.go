package ingestion

type Processor interface {
	Init()
	Start()
	Stop()
}

type RentState struct {
}

type RPOption func(opt *RPOptions)
type RPOptions struct {
	Source Ingestor
}

type RentProcessor struct {
	Source Ingestor
}

func GetRPOptions(opts ...RPOption) *RPOptions {
	rpo := &RPOptions{}

	for _, opt := range opts {
		opt(rpo)
	}

	return rpo
}

func RPWithSource(s Ingestor) RPOption {
	return func(opt *RPOptions) {
		opt.Source = s
	}
}

func AdProcessing(state RentState, data []byte) error {
	return nil
}
