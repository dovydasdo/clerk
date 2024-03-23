package config

type Manager struct {
	Type        string       `yaml:"type"`
	DataSources []DataSource `yaml:"data_sources"`
}
