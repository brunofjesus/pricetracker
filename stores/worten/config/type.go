package config

type ApplicationConfiguration struct {
	MessageQueue struct {
		URL string `yaml:"url"`
	} `yaml:"messagequeue"`
}
