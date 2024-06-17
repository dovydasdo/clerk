package ingestion

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	pad "github.com/dovydasdo/clerk/generated/ad"
	"github.com/dovydasdo/clerk/generated/location"
	queries "github.com/dovydasdo/clerk/sql"
	"github.com/go-co-op/gocron/v2"
	"google.golang.org/protobuf/proto"
)

type RentStateFunc func(ad *pad.Ad, state StateTracker) error
type Action func(data []byte) error

type Processor interface {
	Init() error
	Start() error
	Dump() error
	Stop() error
}

type StateTracker interface {
	Reset()
}

type cityStats struct {
	date          time.Time
	avgPrice      int
	avgPricePerSq float64
	avgFootage    float64
	city          string
	adsCount      uint
	source        string
}

type Message struct {
	source string
	data   []byte
}

type RentState struct {
	statsCity      map[string]*cityStats
	totalProcessed int
}

type RentProcessor struct {
	RcvChan chan Message
	Ctx     context.Context

	// TODO: maybe allow any source to be registered?
	sources []Ingestor
	State   StateTracker
	store   Saver

	stateFuncs    []RentStateFunc
	initStateFunc func(state StateTracker) StateTracker

	l *slog.Logger
	s gocron.Scheduler

	id           string
	dumpInterval string
}

func GetRentProcessor(opts RPOptions) *RentProcessor {
	// TODO: fix this, fails to cast to rent sate if is any
	// state, ok := opts.State.(*RentState)
	// if !ok {
	// 	state = &RentState{
	// 		statsCity:      map[string]*cityStats{},
	// 		totalProcessed: 0,
	// 	}
	// }

	return &RentProcessor{
		store:         opts.Store,
		l:             opts.Logger,
		RcvChan:       make(chan Message),
		Ctx:           opts.Ctx,
		initStateFunc: opts.StateInitFunc,
		stateFuncs:    opts.StateF,
		State:         opts.State,
		sources:       opts.Sources,
		dumpInterval:  opts.DumpInterval,
		s:             opts.Scheduler,
	}
}

func (p RentProcessor) Init() error {
	for _, s := range p.sources {
		p.l.Debug("processor", "message", fmt.Sprintf("init: %+v", s))
		err := s.Init()
		if err != nil {
			return err
		}
	}

	p.State = p.initStateFunc(p.State)

	return nil
}

func (p *RentProcessor) Start() error {
	p.l.Debug("proc", "message", "starting rent processor")
	var err error
	go func() {
		p.l.Debug("proc", "message", "starting to listen for messages")
		for {
			select {
			case msg := <-p.RcvChan:
				p.l.Debug("proc", "received", "message", "rource", msg.source)
				switch msg.source {
				case "rent_ads":
					// parse message to rent data type
					ad := pad.Ad{}
					if err := proto.Unmarshal(msg.data, &ad); err == nil {
						p.l.Debug("proc", "received", "ad", "val", fmt.Sprintf("%v", &ad))
						// update state with current ad
						for _, stateF := range p.stateFuncs {
							err := stateF(&ad, p.State)
							if err != nil {
								break
							}
						}

						// store the parsed ad
						// TODO: consider making saving an async action (submit to some to save queue or smthng)
						err := p.store.Save(&ad)
						if err != nil {
							p.l.Error("proc", "error", err)
							break
						}

						continue
					}
				case "locations":
					loc := location.Location{}
					if err := proto.Unmarshal(msg.data, &loc); err == nil {
						p.l.Debug("proc", "received", "location", "val", fmt.Sprintf("%+v", &loc))
						// Sync locaiton data to db
						err := p.store.Save(&loc)
						if err != nil {
							p.l.Error("proc", "error", err)
							break
						}

						continue

					}
				default:
				}

				p.l.Warn("proc", "message", "received message that could not be cast to any known type")

			case <-p.Ctx.Done():
				for _, s := range p.sources {
					s.Stop()
				}
				break
			}

		}

	}()

	for _, s := range p.sources {
		s.Start(p.RcvChan)
	}

	// TODO: add cron job for dumping sate
	if p.s != nil {
		go func() {
			var def gocron.JobDefinition
			switch p.dumpInterval {
			case "daily":
				// End of day dump
				def = gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(23, 59, 59)))
			case "weekly":
				// End of week dump
				def = gocron.WeeklyJob(1, gocron.NewWeekdays(time.Sunday), gocron.NewAtTimes(gocron.NewAtTime(23, 59, 59)))
			}

			job, err := p.s.NewJob(
				def,
				gocron.NewTask(
					func() {
						p.l.Info("scheduler", "message", fmt.Sprintf("dumping state for: %v", p.id))
						p.Dump()
					},
				),
			)
			if err != nil {
				p.l.Info("proc", "err", err)
			}
			p.l.Info("proc", "message", fmt.Sprintf("started dump job: %v", job.ID()))
			p.s.Start()
		}()
	}

	return err

}

func (p *RentProcessor) Stop() error {
	p.Ctx.Done()
	if p.s != nil {
		p.s.Shutdown()
	}
	return nil
}

func (p RentProcessor) Dump() error {
	// save state and reset it
	err := p.store.Save(p.State)
	if err != nil {
		return err
	}

	p.State.Reset()
	return nil
}

func (s *RentState) Reset() {
	s.totalProcessed = 0
	s.statsCity = map[string]*cityStats{}
}

var SaveAd = func(db *sql.DB, data any, l *slog.Logger) error {
	if ad, ok := data.(*pad.Ad); ok {
		l.Debug("saver", "message", "saving ad")
		res, err := db.Exec(queries.SaveAd, time.Now(), time.Now(), ad.City, ad.Date.AsTime(), ad.Stars, ad.Title, ad.Address, ad.Footage, ad.Rooms, ad.Floor, ad.Specifications, ad.Price, ad.Premium, ad.AdId, ad.Source, ad.Url, ad.BuildingFloors, ad.AdIdUi)
		if res == nil {
			return SaveError{Message: "got nil result from db exec", InnerErr: err}
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return SaveError{Message: "no rows affected", InnerErr: err}
		}

		l.Debug("saver", "message", fmt.Sprintf("rows affected: %v", rows))
		return err
	}
	return errors.New("failed to parse ad")
}

var SaveLocation = func(db *sql.DB, data any, l *slog.Logger) error {
	if loc, ok := data.(*location.Location); ok {
		_, err := db.Exec(queries.SaveLocation, loc.Id, time.Now(), time.Now(), loc.Lat, loc.Lng)
		if err != nil {
			return SaveError{Message: "got nil result from db exec", InnerErr: err}
		}
		return nil
	}
	return errors.New("failed to parse location")
}

var SaveRentState = func(db *sql.DB, state any, l *slog.Logger) error {
	l.Info("saver", "message", "saving state")
	if st, ok := state.(*RentState); ok {
		l.Info("saver", "message", fmt.Sprintf("dumping state, processed ads: %v", st.totalProcessed))
		valueStrings := make([]string, 0, len(st.statsCity))
		valueArgs := make([]interface{}, 0, len(st.statsCity)*7)
		for _, city := range st.statsCity {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")
			valueArgs = append(valueArgs, time.Now())
			valueArgs = append(valueArgs, city.avgPrice)
			valueArgs = append(valueArgs, city.avgPricePerSq)
			valueArgs = append(valueArgs, city.avgFootage)
			valueArgs = append(valueArgs, city.city)
			valueArgs = append(valueArgs, city.date)
			valueArgs = append(valueArgs, city.adsCount)
		}

		_, err := db.Exec(fmt.Sprintf(queries.SaveRentState, strings.Join(valueStrings, ",")), valueArgs...)
		if err != nil {
			return SaveError{Message: "failed to save rent sate", InnerErr: err}
		}

		return nil
	}

	return errors.New("failed to parse rent state")
}

var InitSate = func(state StateTracker) StateTracker {
	return &RentState{
		statsCity:      make(map[string]*cityStats),
		totalProcessed: 0,
	}
}

var UpdateRentState = func(ad *pad.Ad, state StateTracker) error {
	rState, ok := state.(*RentState)
	if !ok {
		return errors.New("failed to cast to rent state")
	}

	cStats, ok := rState.statsCity[ad.City]
	if !ok {
		date := time.Now()
		date = time.Date(
			date.Year(), date.Month(), date.Day(),
			0, 0, 0, 0, date.Location(),
		)

		cs := cityStats{
			date:          date,
			avgPrice:      int(ad.Price),
			avgPricePerSq: float64(ad.Price) / ad.Footage,
			avgFootage:    ad.Footage,
			city:          ad.City,
			adsCount:      1,
			source:        ad.Source,
		}

		rState.statsCity[ad.City] = &cs
		rState.totalProcessed++
		return nil
	}

	// This might be dumb
	cStats.avgPrice = (cStats.avgPrice*int(cStats.adsCount) + int(ad.Price)) / int(cStats.adsCount+1)
	cStats.avgFootage = (cStats.avgFootage*float64(cStats.adsCount) + float64(ad.Footage)) / float64(cStats.adsCount+1)
	cStats.avgPricePerSq = float64(cStats.avgPrice) / cStats.avgFootage
	cStats.adsCount++
	rState.totalProcessed++
	return nil
}
