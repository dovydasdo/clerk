package ingestion

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// TODO: make this more configurable

type NatsKVSync struct {
	nc *nats.Conn
	kv jetstream.KeyValue

	ctx     context.Context
	watcher jetstream.KeyWatcher
	bName   string
	name    string
	url     string
	subject string
}

func GetNatsKVSync(options *NKVCacheOptions) *NatsKVSync {
	return &NatsKVSync{
		ctx:     options.Ctx,
		name:    options.Name,
		url:     options.Url,
		subject: options.Subject,
	}
}

func (c *NatsKVSync) Init() error {
	// TODO: maybe use a single conneciton for all NATS sources?
	nc, err := nats.Connect(c.url)
	if err != nil {
		return err
	}
	c.nc = nc
	js, err := jetstream.New(nc)
	if err != nil {
		return err
	}
	kv, err := js.CreateOrUpdateKeyValue(c.ctx, jetstream.KeyValueConfig{
		Bucket: c.name,
	})
	if err != nil {
		return err
	}
	c.kv = kv
	w, err := c.kv.Watch(c.ctx, c.subject, jetstream.UpdatesOnly())
	if err != nil {
		return err
	}

	c.watcher = w
	return nil
}

func (c *NatsKVSync) Start(sChan chan Message) error {
	var err error
	go func() {
		for {
			select {
			case kvEntry := <-c.watcher.Updates():
				if kvEntry == nil {
					fmt.Printf("recieved: nil \n")
					continue
				}
				fmt.Printf("recieved: %+v \n", kvEntry.Value())
				sChan <- Message{source: c.name, data: kvEntry.Value()}
			case <-c.ctx.Done():
				return
			}
		}
	}()

	return err
}

func (c *NatsKVSync) Stop() error {
	c.ctx.Done()
	c.nc.Drain()
	return c.watcher.Stop()
}
