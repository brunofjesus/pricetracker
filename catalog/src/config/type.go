package config

type ApplicationConfiguration struct {
	Database struct {
		DSN string `yaml:"DSN"`
	} `yaml:"database"`
}
