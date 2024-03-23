package management

import "testing"

type Processor interface {
	Init() error
	Start() error
	Dump() error
	Stop() error
}

type MockRentProcessor struct {
}

func (p MockRentProcessor) Start() error {
	return nil
}

func (p MockRentProcessor) Init() error {
	return nil
}

func (p MockRentProcessor) Dump() error {
	return nil
}

func (p MockRentProcessor) Stop() error {
	return nil
}

func TestRentManager(t *testing.T) {
	proc := MockRentProcessor{}

	mngrOpts := GetRMOptions(
		RMWithProcessor(
			proc,
		),
	)

	manager := GetRentManager(*mngrOpts)
	err := manager.Init()
	if err != nil {
		t.Errorf("manager failed at init")
	}

	err = manager.Start()
	if err != nil {
		t.Errorf("manager failed at start")
	}

	err = manager.Stop()
	if err != nil {
		t.Errorf("manager failed at stop")
	}
}
