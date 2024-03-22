package ingestion

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

// TODO: make this more configurable

type NatsKVSync struct {
	kv      jetstream.KeyValue
	ctx     context.Context
	watcher jetstream.KeyWatcher
	bName   string
	name    string
}

func GenNatsKVSync(kv jetstream.KeyValue, ctx context.Context, name string) *NatsKVSync {
	return &NatsKVSync{
		kv:   kv,
		ctx:  ctx,
		name: name,
	}
}

func (c *NatsKVSync) Init() error {
	w, err := c.kv.Watch(c.ctx, c.bName)
	c.watcher = w
	return err
}

func (c NatsKVSync) Start(sChan chan Message) error {
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

func (c NatsKVSync) Stop() error {
	c.ctx.Done()
	return c.watcher.Stop()
}
