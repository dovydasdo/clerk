package ingestion

import (
	"context"

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
}

func GetNatsKVSync(options *NKVCacheOptions) *NatsKVSync {
	return &NatsKVSync{
		ctx:  options.Ctx,
		name: options.Name,
		url:  options.Url,
	}
}

func (c *NatsKVSync) Init() error {
	nc, err := nats.Connect(c.url)
	if err != nil {
		return err
	}
	c.nc = nc

	js, _ := jetstream.New(nc)
	kv, err := js.CreateOrUpdateKeyValue(c.ctx, jetstream.KeyValueConfig{
		Bucket: c.name,
	})
	if err != nil {

		return err
	}

	c.kv = kv

	w, err := c.kv.Watch(c.ctx, c.bName)
	c.watcher = w
	return err
}

func (c *NatsKVSync) Start(sChan chan Message) error {
	var err error
	go func() {
		select {
		case kvEntry := <-c.watcher.Updates():
			sChan <- Message{source: c.name, data: kvEntry.Value()}
		case <-c.ctx.Done():
			return
		}
	}()

	return err
}

func (c *NatsKVSync) Stop() error {
	c.ctx.Done()
	c.nc.Drain()
	return c.watcher.Stop()
}
