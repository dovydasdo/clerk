package ingestion

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	pad "github.com/dovydasdo/clerk/generated/ad"
	"github.com/dovydasdo/clerk/generated/location"
	queries "github.com/dovydasdo/clerk/sql"
	"google.golang.org/protobuf/proto"
)

type RentStateFunc func(ad *pad.Ad, state any) error
type DumpFunc func(state any) error
type Action func(data []byte) error

type Processor interface {
	Init() error
	Start() error
	Dump() error
	Stop() error
}

type cityStats struct {
	date          time.Time
	avgPrice      int
	avgPricePerSq float64
	avgFootage    float64
	city          string
	createdAt     time.Time
	adsCount      uint
	source        string
}

type RentState struct {
	statsCity cityStats
}

type RentProcessor struct {
	RcvChan chan []byte
	Ctx     context.Context

	// TODO: maybe allow any source to be registered?
	source  Ingestor
	kvStore CacheSyncer
	state   any
	store   Saver

	stateFuncs    []RentStateFunc
	initStateFunc func(state any)

	l *slog.Logger
}

func GetRentProcessor(opts RPOptions) *RentProcessor {
	return &RentProcessor{
		source:        opts.Source,
		store:         opts.Store,
		l:             opts.Logger,
		RcvChan:       make(chan []byte),
		Ctx:           opts.Ctx,
		initStateFunc: opts.StateInitFunc,
		stateFuncs:    opts.StateF,
		state:         opts.State,
		kvStore:       opts.KVStore,
	}
}

func (p *RentProcessor) Start() error {
	var err error
	go func() {
		for {
			select {
			case msg := <-p.RcvChan:
				p.l.Debug("proc", "received", msg)
				// parse message to rent data type
				ad := pad.Ad{}
				if err := proto.Unmarshal(msg, &pad.Ad{}); err != nil {
					p.l.Debug("proc", "received", "ad")
					// update state with current ad
					for _, stateF := range p.stateFuncs {
						err := stateF(&ad, p.state)
						if err != nil {
							break
						}
					}

					// store the parsed ad
					// TODO: consider making saving an async action (submit to some to save queue or smthng)
					err := p.store.Save(&ad)
					if err != nil {
						break
					}

					continue
				}

				loc := location.Location{}
				if err := proto.Unmarshal(msg, &loc); err != nil {
					p.l.Debug("proc", "received", "location")
					// Sync locaiton data to db
					err := p.store.Save(&ad)
					if err != nil {
						break
					}

					continue

				}

				p.l.Warn("proc", "message", "received message that could not be cast to any known type")
			case <-p.Ctx.Done():
				p.source.Stop()
				p.kvStore.Stop()
				break
			}

		}

	}()
	p.source.Init()
	p.source.Start(p.RcvChan)
	p.kvStore.Start(p.RcvChan)

	p.l.Debug("proc", "message", "starting to listen for messages")
	fmt.Printf("starting to listen for messages \n")

	return err

}

func (p *RentProcessor) Stop() error {
	p.Ctx.Done()
	return nil
}

func (p RentProcessor) Dump() error {
	// save state and reset it
	err := p.store.Save(p.state)
	if err != nil {
		return err
	}

	p.initStateFunc(p.state)
	return nil
}

var SaveAd = func(db *sql.DB, data any) error {
	if ad, ok := data.(*pad.Ad); ok {
		_, err := db.Exec(queries.SaveAd, time.Now(), time.Now(), ad.City, ad.Date, ad.Stars, ad.Title, ad.Address, ad.Footage, ad.Rooms, ad.Floor, ad.Specifications, ad.Price, ad.Premium, ad.AdId, ad.Source, ad.Url, ad.BuildingFloors, ad.AdIdUi)
		return err
	}
	return errors.New("failed to parse ad")
}

var SaveLocation = func(db *sql.DB, data any) error {
	if loc, ok := data.(*location.Location); ok {
		_, err := db.Exec(queries.SaveLocation, loc.Id, time.Now(), time.Now(), loc.Lat, loc.Lng)
		return err
	}
	return errors.New("failed to parse location")
}

var SaveState = func(db *sql.DB, state any) error {
	// save state

	return nil
}

var InitSate = func(state any) {
	// cast state to rent state and make a new instance
}

var UpdateRentState = func(ad *pad.Ad, state any) error {
	return nil
}
