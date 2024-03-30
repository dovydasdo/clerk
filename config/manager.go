package config

type Manager struct {
	Type       string      `mapstructure:"type"`
	Processors []Processor `mapstructure:"processors"`
}
