package config

type DataSource struct {
	Name   string `mapstructure:"name"`
	Stream string `mapstructure:"stream"`
	KVName string `mapstructure:"kv_name"`
	Store  string `mapstructure:"store"`
	Type   string `mapstructure:"type"`
	Dump   string `mapstructure:"dump"`
}
