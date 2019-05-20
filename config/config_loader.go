package config

import (
	"gitlab.com/gps2.0/errcheck"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

var appConfigInstance *AppConfig
var once sync.Once

// GetAppConfig returns an AppConfig instance
func GetAppConfig() *AppConfig {
	once.Do(func() {
		loadConfig()
	})

	return appConfigInstance
}

func loadConfig() {
	yamlFile, err := ioutil.ReadFile("config.yml")

	errcheck.CheckError(err)

	err = yaml.Unmarshal(yamlFile, &appConfigInstance)
	errcheck.CheckError(err)
}
