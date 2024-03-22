package ingestion

import (
	"context"
	"database/sql"

	"github.com/dovydasdo/clerk/generated/location"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
)

// TODO: make this more configurable

type CacheSyncer interface {
	Start(q string) error
	Stop() error
	Sync() error
}

type NatsKVSync struct {
	db      *sql.DB
	kv      jetstream.KeyValue
	ctx     context.Context
	watcher jetstream.KeyWatcher
}

func GenNatsKVSync(db *sql.DB, kv jetstream.KeyValue, ctx context.Context) *NatsKVSync {
	return &NatsKVSync{
		db:  db,
		kv:  kv,
		ctx: ctx,
	}
}

func (c NatsKVSync) Start(query string) error {
	w, _ := c.kv.Watch(c.ctx, query)
	c.watcher = w
	var err error
	go func() {
		select {
		case kvEntry := <-c.watcher.Updates():
			loc := location.Location{}
			err = proto.Unmarshal(kvEntry.Value(), &loc)
			if err != nil {
				// TODO: handle more gracefully
				return
			}

			//check if already in db
			//update/create new entry
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
