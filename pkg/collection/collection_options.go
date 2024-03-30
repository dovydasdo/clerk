package collection

import (
	"github.com/dovydasdo/clerk/config"
)

type CollectionOption func(opts *CollectionOptions)
type CollectionOptions struct {
}

func GetCollectionOptions(cfg config.CollectionConfig) *CollectionOptions {
	opts := &CollectionOptions{}

	return opts
}
