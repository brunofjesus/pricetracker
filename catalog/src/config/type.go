package config

type ApplicationConfiguration struct {
	Database struct {
		DSN      string `yaml:"dsn"`
		Attempts int    `yaml:"attempts"`
	} `yaml:"database"`
}
