package config

type ApplicationConfiguration struct {
	MessageQueue struct {
		URL string `yaml:"url"`
	} `yaml:"messagequeue"`
	CategoriesSitemap   string `yaml:"categoriesSitemap"`
	PolitenessDelayMs   int64  `yaml:"politenessDelayMs"`
	LoopIntervalMinutes int64  `yaml:"loopIntervalMinutes"`
}
