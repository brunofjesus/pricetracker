package config

type ApplicationConfiguration struct {
	MessageQueue struct {
		URL string `yaml:"url"`
	} `yaml:"messagequeue"`
	Store struct {
		Department string `yaml:"department"`
		Country    string `yaml:"country"`
		Slug       string `yaml:"slug"`
		Name       string `yaml:"name"`
	} `yaml:"store"`
	UserAgent           string `yaml:"userAgent"`
	LoopIntervalMinutes int64  `yaml:"loopIntervalMinutes"`
	PolitenessDelayMs   int64  `yaml:"politenessDelayMs"`
}
