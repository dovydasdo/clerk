package ingestion

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Procedure func(state any, data []byte) error
type DumpF func(state any) error

type Ingestor interface {
	Init() error
	Start() error
	Stop() error
	// Register(proc []Procedure) error <--- maybe later...
	Dump() error
}

// TODO: some more configurations may be needed
type NATSIngestor struct {
	// State to track stats for all data
	// To be dumped on Stop() or Dump()
	State any

	// NATS configuration
	ServerUrl  string
	StreamName string
	Subject    string

	// Processing functions
	DFunc DumpF
	PFunc Procedure

	conn *nats.Conn
	js   jetstream.JetStream
	cons jetstream.Consumer
	ctx  context.Context
	l    *slog.Logger

	// stop functions
	consStop func()
}

func GetNATSIngestor(opts NIOptions) *NATSIngestor {
	return &NATSIngestor{
		ServerUrl:  opts.ServerUrl,
		StreamName: opts.StreamName,
		Subject:    opts.Subject,
		DFunc:      opts.DFunc,
		PFunc:      opts.PFunc,
		ctx:        opts.Ctx,
	}
}

func (n *NATSIngestor) Init() error {
	nc, err := nats.Connect(n.ServerUrl)
	if err != nil {
		return err
	}
	n.conn = nc
	js, err := jetstream.New(nc)
	if err != nil {
		return err
	}
	n.js = js
	cons, err := js.CreateOrUpdateConsumer(n.ctx, n.StreamName, jetstream.ConsumerConfig{})
	if err != nil {
		return err
	}

	n.cons = cons

	return nil
}

func (n *NATSIngestor) Start() error {
	// Start ingesting data from stream
	cc, _ := n.cons.Consume(func(msg jetstream.Msg) {
		n.l.Debug("nats", "ingested from", msg.Subject())

		err := n.PFunc(n.State, msg.Data())
		if err != nil {
			n.l.Error("nats.process", "error", err.Error())
		}
	})

	n.consStop = cc.Stop

	return nil
}

func (n NATSIngestor) Stop() error {
	n.consStop()
	n.conn.Drain()
	n.ctx.Done()

	return nil
}

func (n *NATSIngestor) Dump() error {
	err := n.DFunc(n.State)

	// This may be a bad idea
	n.State = nil

	return err

}
