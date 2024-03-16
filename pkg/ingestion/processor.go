package ingestion

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	pad "github.com/dovydasdo/clerk/generated/ad"
	queries "github.com/dovydasdo/clerk/sql"
	"google.golang.org/protobuf/proto"
)

type Action func(data []byte) error

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
	Store  Saver
}

type RentProcessor struct {
	RcvChan chan []byte
	Ctx     context.Context

	source Ingestor
	store  Saver
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

func RPWithSaver(s Saver) RPOption {
	return func(opt *RPOptions) {
		opt.Store = s
	}
}

func GetRentProcessor(opts RPOptions) *RentProcessor {
	return &RentProcessor{
		source: opts.Source,
		store:  opts.Store,
	}
}

func (p *RentProcessor) Start() error {
	p.source.Init()
	p.source.Start(p.RcvChan)

	for {
		select {
		case msg := <-p.RcvChan:
			ad := pad.Ad{}
			if err := proto.Unmarshal(msg, &ad); err != nil {
				//temp
				fmt.Printf("failed to unmarshall message")
				continue
			}

			// TODO: consider making saving an async action (submit to some to save queue or smthng)
			err := p.store.Save(&ad)
			if err != nil {
				return err

			}
		}

	}
}

func GetProcedureWithActions(acts []Action) Procedure {
	return func(state any, data []byte) error {
		for _, action := range acts {
			err := action(data)
			if err != nil {
				break
			}
		}

		return nil
	}
}

func GetSaveAction(s Saver) Action {
	return func(data []byte) error {

		s.Save(data)
		return nil
	}
}

func AdProcessing(state RentState, data []byte) error {
	return nil
}

var SaveAd = func(db *sql.DB, data any) error {
	if ad, ok := data.(*pad.Ad); ok {
		_, err := db.Exec(queries.SaveAd, time.Now(), time.Now(), ad.City, ad.Date, ad.Stars, ad.Title, ad.Address, ad.Footage, ad.Rooms, ad.Floor, ad.Specifications, ad.Price, ad.Premium, ad.AdId, ad.Source, ad.Url, ad.BuildingFloors, ad.AdIdUi)
		return err
	}

	return errors.New("failed to parse ad")
}
