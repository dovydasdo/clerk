package ingestion

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

// TODO: make this more configurable

type CacheSyncer interface {
	Start(c chan []byte) error
	Stop() error
	Sync() error
}

type NatsKVSync struct {
	kv      jetstream.KeyValue
	ctx     context.Context
	watcher jetstream.KeyWatcher
	bName   string
}

func GenNatsKVSync(kv jetstream.KeyValue, ctx context.Context) *NatsKVSync {
	return &NatsKVSync{
		kv:  kv,
		ctx: ctx,
	}
}

func (c NatsKVSync) Start(sChan chan []byte) error {
	w, _ := c.kv.Watch(c.ctx, c.bName)
	c.watcher = w
	var err error
	go func() {
		select {
		case kvEntry := <-c.watcher.Updates():
			sChan <- kvEntry.Value()
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

func (c NatsKVSync) Sync() error {
	// TODO: implement manual sync
	return nil
}
