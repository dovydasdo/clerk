package ingestion

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/dovydasdo/clerk/generated/ad"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			fmt.Printf("sending: %v \n", val)
			if sChan == nil {
				fmt.Printf("no chan \n")
			}
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
	saveF func(data any)
}

func (s MockSaver) Save(data any) error {
	s.saveF(data)
	return nil
}

func (s MockSaver) Close() error {
	return nil
}

var testIngestor *MockIngestor
var testProc *RentProcessor

type TestRentState struct {
	average int
	count   int
}

func TestIngestion(t *testing.T) {
	c, cf := context.WithCancel(context.Background())
	testIngestor = &MockIngestor{
		ctx:     c,
		cancelF: cf,
	}

	var priceTotal int32
	adsNum := 10

	for i := 0; i < adsNum; i++ {
		price := int32(i * 10)
		testAd := ad.Ad{
			Id:             uint32(i),
			CreatedAt:      &timestamppb.Timestamp{},
			UpdatedAt:      &timestamppb.Timestamp{},
			DeletedAt:      &timestamppb.Timestamp{},
			City:           "",
			Date:           &timestamppb.Timestamp{},
			Stars:          0,
			Title:          "",
			Address:        "",
			Footage:        0,
			Rooms:          0,
			Floor:          0,
			BuildingFloors: 0,
			Specifications: "",
			Price:          price,
			Premium:        false,
			AdId:           "",
			AdIdUi:         0,
			Source:         "",
			Url:            "",
		}

		priceTotal = priceTotal + price
		data, _ := proto.Marshal(&testAd)
		sendData = append(sendData, data)
	}

	averagePrice := priceTotal / int32(adsNum)

	received := 0
	procOpts := GetRPOptions(
		RPWithSource(testIngestor),
		RPWithSaver(MockSaver{
			saveF: func(data any) {
				if state, ok := data.(*TestRentState); ok {
					if state.average != int(averagePrice) {
						t.Errorf("average incorrect, wanted: %v, got: %v", averagePrice, state.average)
					}

				}

				if _, ok := data.(*ad.Ad); ok {
					received++
				}
			},
		}),
		RPWithState(&TestRentState{}),
		RPWithStateInitF(func(s any) {
			s = &TestRentState{}
		}),
		RPWithStateF(func(ad *ad.Ad, state any) error {
			if state, ok := state.(*TestRentState); ok {
				currSum := state.average * state.count
				state.count++
				state.average = (currSum + int(ad.Price)) / state.count
			}
			return nil
		}),
		RPWithLogger(slog.Default()),
		RPWithCtx(context.Background()),
	)

	testProc = GetRentProcessor(*procOpts)
	err := testProc.Start()
	if err != nil {
		t.Errorf("processor failed: %v", err)
	}

	err = testProc.Dump()
	if err != nil {
		t.Errorf("processor failed at dump: %v", err)
	}

	if received != adsNum {
		t.Errorf("did not get all ads, wanted: %v, got: %v", adsNum, received)
	}
}
