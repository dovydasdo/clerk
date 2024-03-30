package config

type DataSource struct {
	Name       string `mapstructure:"name"`
	Stream     string `mapstructure:"stream"`
	Subject    string `mapstructure:"subject"`
	Type       string `mapstructure:"type"`
	Dump       string `mapstructure:"dump"`
	SourceType string `mapstructure:"source_type"`
}
