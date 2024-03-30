package ingestion

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Ingestor interface {
	Init() error
	Start(sChan chan Message) error
	Stop() error
}

// TODO:	some more configurations may be needed
//		listen to kv store and sync it with db (maybe seperate struct for that)
type NATSIngestor struct {
	ServerUrl  string
	StreamName string
	Subject    string

	conn      *nats.Conn
	js        jetstream.JetStream
	cons      jetstream.Consumer
	ctx       context.Context
	l         *slog.Logger
	name      string
	batchSize int
}

func GetNATSIngestor(opts *NIOptions) *NATSIngestor {
	return &NATSIngestor{
		ServerUrl:  opts.ServerUrl,
		StreamName: opts.StreamName,
		Subject:    opts.Subject,
		ctx:        opts.Ctx,
		l:          opts.Logger,
		name:       opts.Name,
		batchSize:  50,
	}
}

func (n *NATSIngestor) Init() error {
	n.l.Debug("ingestor", "message", fmt.Sprintf("starting nats connection at: %v", n.ServerUrl))
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
	n.l.Debug("ingestor", "message", fmt.Sprintf("starting stream consumer for: %v", n.StreamName))
	n.l.Debug("ingestor", "message", fmt.Sprintf("with subject: %v", n.Subject))
	cons, err := js.CreateOrUpdateConsumer(n.ctx, n.StreamName, jetstream.ConsumerConfig{
		FilterSubject: n.Subject,
		Durable:       fmt.Sprintf("%s_consumer", n.StreamName),
	})

	if err != nil {
		return err
	}

	n.cons = cons

	n.l.Debug("ingestor", "message", fmt.Sprintf("ingestor inited: %+v", n))
	return nil
}

func (n *NATSIngestor) Start(sChan chan Message) error {
	// Start ingesting data from stream
	go func() {
		n.l.Debug("ingestor", "message", "starting nats consumer")
		for {
			select {
			case <-n.ctx.Done():
				close(sChan)
				break
			default:
				n.l.Debug("ingetsor", "message", "fetching")
				msgs, _ := n.cons.Fetch(n.batchSize)
				for msg := range msgs.Messages() {
					if msg != nil {
						n.l.Debug("ingestor", "ingested from", msg.Subject(), "source name", n.name)
						sChan <- Message{source: n.name, data: msg.Data()}
						// Using DoubleAck to prevent duplicate data
						err := msg.DoubleAck(n.ctx)
						if err != nil {
							n.l.Warn("ingetsor", "message", "failed to ack")
						}

						continue
					}

					n.l.Warn("ingetsor", "message", "received nil message")
				}
			}
		}
	}()

	return nil
}

func (n NATSIngestor) Stop() error {
	n.conn.Drain()
	n.ctx.Done()

	return nil
}
