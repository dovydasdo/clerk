package management

import (
	i "github.com/dovydasdo/clerk/pkg/ingestion"
)

// Contains all of the processors for rent data sites
type RentManager struct {
	Processors []i.Processor
}

func GetRentManager() *RentManager {
	rm := &RentManager{
		Processors: make([]i.Processor, 0),
	}

	processorOpts := i.GetRPOptions()
	processor := i.GetRentProcessor(*processorOpts)

	rm.Processors = append(rm.Processors, processor)
	return rm
}

func (rm RentManager) Init() error {
	return nil
}

func (rm RentManager) Start() error {
	return nil
}

func (rm RentManager) Stop() error {
	return nil
}
