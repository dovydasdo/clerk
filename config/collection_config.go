package config

type CollectionConfig struct {
	Managers []Manager `mapstructure:"managers"`
}
