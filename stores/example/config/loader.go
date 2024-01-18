package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"sync"
)

var once sync.Once
var instance *ApplicationConfiguration

func GetApplicationConfiguration() *ApplicationConfiguration {
	once.Do(func() {
		var err error
		instance, err = loadConfig()

		if err != nil {
			panic(err)
		}
	})

	return instance
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
