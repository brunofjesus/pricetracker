package config

type ApplicationConfiguration struct {
	MessageQueue struct {
		URL string `yaml:"url"`
	} `yaml:"messagequeue"`
	LoopIntervalMs int64 `yaml:"loopIntervalMs"`
}
