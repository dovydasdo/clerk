package ingestion

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Procedure func(state any, data []byte) error

type Ingestor interface {
	Init() error
	Start(sChan chan []byte) error
	Stop() error
}

// TODO:	some more configurations may be needed
//		listen to kv store and sync it with db (maybe seperate struct for that)
type NATSIngestor struct {
	ServerUrl  string
	StreamName string
	Subject    string

	conn *nats.Conn
	js   jetstream.JetStream
	cons jetstream.Consumer
	ctx  context.Context
	l    *slog.Logger

	consStop func()
}

func GetNATSIngestor(opts NIOptions) *NATSIngestor {
	return &NATSIngestor{
		ServerUrl:  opts.ServerUrl,
		StreamName: opts.StreamName,
		Subject:    opts.Subject,
		ctx:        opts.Ctx,
		l:          opts.Logger,
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

func (n *NATSIngestor) Start(sChan chan []byte) error {
	// Start ingesting data from stream

	go func() {
		cc, _ := n.cons.Consume(func(msg jetstream.Msg) {
			n.l.Debug("nats", "ingested from", msg.Subject())
			// Read unitl ctx is stopped
			select {
			case <-n.ctx.Done():
				close(sChan)
				break
			default:
				sChan <- msg.Data()
			}
		})

		n.consStop = cc.Stop
	}()

	return nil
}

func (n NATSIngestor) Stop() error {
	n.consStop()
	n.conn.Drain()
	n.ctx.Done()

	return nil
}
