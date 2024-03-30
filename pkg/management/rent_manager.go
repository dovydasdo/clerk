package management

import (
	"context"
	"log/slog"

	i "github.com/dovydasdo/clerk/pkg/ingestion"
)

// Starts and manages processors
type RentManager struct {
	Processors []i.Processor
	ctx        context.Context
	l          *slog.Logger
}

func GetRentManager(opts RMOptions) *RentManager {
	rm := &RentManager{
		Processors: opts.Processors,
		ctx:        opts.Ctx,
		l:          opts.Logger,
	}

	return rm
}

func (rm RentManager) Init() error {
	for _, p := range rm.Processors {
		err := p.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

func (rm RentManager) Start() error {
	for _, p := range rm.Processors {
		err := p.Start()
		if err != nil {
			return err
		}
	}

	// TODO: have so loop to check for changes in config/some state api

	return nil
}

func (rm RentManager) Stop() error {
	return nil
}
