package ingestion

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/dovydasdo/clerk/generated/ad"
	"github.com/dovydasdo/clerk/generated/location"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type leveler struct{}

func (l leveler) Level() slog.Level {
	return -4
}

var sendData [][]byte
var sendLocationData [][]byte

type MockIngestor struct {
	ctx     context.Context
	cancelF context.CancelFunc
	data    [][]byte
	name    string
}

func (i *MockIngestor) Init() error {
	return nil
}

func (i *MockIngestor) Start(sChan chan Message) error {
	go func() {
		for _, val := range i.data {
			select {
			case <-i.ctx.Done():
				break
			default:
				fmt.Printf("sending: %v \n", val)
				if sChan == nil {
					fmt.Printf("no chan \n")
				}

				sChan <- Message{source: i.name, data: val}
			}
		}
	}()

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
	logFile := os.Stdout
	logger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: leveler{}}))
	c, cf := context.WithCancel(context.Background())

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

	locNum := 10

	for i := 0; i < locNum; i++ {
		testLoc := location.Location{
			Id:        int32(i),
			CreateAt:  &timestamppb.Timestamp{},
			UpdatedAt: &timestamppb.Timestamp{},
			DeletedAt: &timestamppb.Timestamp{},
			Lat:       0,
			Lng:       0,
		}

		data, _ := proto.Marshal(&testLoc)
		sendLocationData = append(sendLocationData, data)
	}

	testIngestor = &MockIngestor{
		ctx:     c,
		cancelF: cf,
		data:    sendData,
		name:    "rent_ads",
	}

	testLocIngestor := &MockIngestor{
		ctx:     c,
		cancelF: cf,
		data:    sendLocationData,
		name:    "locations",
	}

	fmt.Printf("testLocIngestor: %v\n", testLocIngestor)

	received := 0
	receivedLoc := 0
	procOpts := GetRPOptions(
		RPWithSource(testIngestor),
		RPWithSource(testLocIngestor),
		RPWithSaver(MockSaver{
			saveF: func(data any) {
				if state, ok := data.(*TestRentState); ok {
					if state.average != int(averagePrice) {
						t.Errorf("average incorrect, wanted: %v, got: %v", averagePrice, state.average)
					}
				}

				if ad, ok := data.(*ad.Ad); ok {
					logger.Debug("test", "message", "saving", "ad", fmt.Sprintf("%v", ad))
					received++
				}

				if _, ok := data.(*location.Location); ok {
					receivedLoc++
				}
			},
		}),
		RPWithState(&TestRentState{}),
		RPWithStateInitF(func(s any) {
			s = &TestRentState{}
		}),
		RPWithStateF(func(ad *ad.Ad, state any) error {
			if state, ok := state.(*TestRentState); ok {
				logger.Debug("test", "message", "processing state", "state", fmt.Sprintf("%v", state))
				logger.Debug("test", "message", "processing state", "ad", fmt.Sprintf("%v", ad))
				currSum := state.average * state.count
				state.count++
				state.average = (currSum + int(ad.Price)) / state.count
			}
			return nil
		}),
		RPWithLogger(logger),
		RPWithCtx(context.Background()),
	)

	testProc = GetRentProcessor(*procOpts)
	err := testProc.Start()
	if err != nil {
		t.Errorf("processor failed: %v", err)
	}

	// TODO: maybe have some drain funcitonality
	time.Sleep(time.Second)

	err = testProc.Dump()
	if err != nil {
		t.Errorf("processor failed at dump: %v", err)
	}

	if received != adsNum {
		t.Errorf("did not get all ads, wanted: %v, got: %v", adsNum, received)
	}

	if receivedLoc != locNum {
		t.Errorf("did not get all locations, wanted: %v, got: %v", locNum, receivedLoc)
	}
}
