package config

type ApplicationConfiguration struct {
	MessageQueue struct {
		URL string `yaml:"url"`
	} `yaml:"messagequeue"`
	LoopIntervalMinutes int64 `yaml:"loopIntervalMinutes"`
	PolitenessDelay     struct {
		PageMs     int64 `yaml:"pageMs"`
		CategoryMs int64 `yaml:"categoryMs"`
	} `yaml:"politenessDelay"`
}
