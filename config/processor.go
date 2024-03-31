package config

type Processor struct {
	Stores      []Store      `mapstructure:"stores"`
	DataSources []DataSource `mapstructure:"data_sources"`
	Id          string       `mapstructure:"id"`
	StateDump   string       `mapstructure:"state_dump"`
}
