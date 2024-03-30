package config

type Store struct {
	Type     string `mapstructure:"type"`
	Endpoint string `mapstructure:"endpoint"`
}
