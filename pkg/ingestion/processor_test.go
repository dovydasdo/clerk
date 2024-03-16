package ingestion

import (
	"context"
)

var sendData [][]byte

type MockIngestor struct {
	ctx     context.Context
	cancelF context.CancelFunc
}

func (i *MockIngestor) Init() error {
	return nil
}

func (i *MockIngestor) Start(sChan chan []byte) error {
	for _, val := range sendData {
		select {
		case <-i.ctx.Done():
			break
		default:
			sChan <- val
		}
	}

	return nil
}

func (i *MockIngestor) Stop() error {
	i.cancelF()
	return nil
}

type MockSaver struct {
}

func (s MockSaver) Save(data any) error {
	return nil
}

func (s MockSaver) Close() error {
	return nil
}

var testIngestor *MockIngestor
var testProc *RentProcessor

func init() {
	c, cf := context.WithCancel(context.Background())
	testIngestor = &MockIngestor{
		ctx:     c,
		cancelF: cf,
	}

	procOpts := GetRPOptions(
		RPWithSource(testIngestor),
		RPWithSaver(MockSaver{}),
	)

	testProc = GetRentProcessor(*procOpts)

	sendData = [][]byte{[]byte("message one"), []byte("message two"), []byte("message three")}
}

func TestIngestion() {

}
