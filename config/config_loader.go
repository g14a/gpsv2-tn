package config

import (
	"gitlab.com/gpsv2/errcheck"
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

// loadConfig loads the data in the yaml file into a struct
// returns the app instance once if it is ready.
func loadConfig() {
	yamlFile, err := ioutil.ReadFile("config.yml")

	errcheck.CheckError(err)

	err = yaml.Unmarshal(yamlFile, &appConfigInstance)

	errcheck.CheckError(err)
}
