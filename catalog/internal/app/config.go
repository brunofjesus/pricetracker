package app

import (
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

var configOnce sync.Once
var configInstance *ApplicationConfiguration

type ApplicationConfiguration struct {
	Database struct {
		DSN      string `yaml:"dsn"`
		Attempts int    `yaml:"attempts"`
		Migrate  bool   `yaml:"migrate"`
	} `yaml:"database"`
	MessageQueue struct {
		URL         string `yaml:"url"`
		ManualAck   bool   `yaml:"manualAck"`
		ThreadCount int    `yaml:"threadCount"`
	} `yaml:"messagequeue"`
}

func GetApplicationConfiguration() *ApplicationConfiguration {
	configOnce.Do(func() {
		var err error
		configInstance, err = loadConfig()

		if err != nil {
			panic(err)
		}
	})

	return configInstance
}

func loadConfig() (*ApplicationConfiguration, error) {
	var config ApplicationConfiguration

	contents, err := os.ReadFile("config.yaml")

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(contents, &config)

	return &config, err
}
