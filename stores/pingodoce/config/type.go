package config

type ApplicationConfiguration struct {
	MessageQueue struct {
		URL string `yaml:"url"`
	} `yaml:"messagequeue"`
	LoopIntervalMs  int64 `yaml:"loopIntervalMs"`
	PolitenessDelay struct {
		PageMs     int64 `yaml:"pageMs"`
		CategoryMs int64 `yaml:"categoryMs"`
	} `yaml:"politenessDelay"`
}
