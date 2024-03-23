package config

type Collection struct {
	DataSources []DataSource `mapstructure:"data_sources"`
}
