// Package config is the run time representation of
// the config.yml file and takes in all the details
// of the databases and tcp server and port.
// DO NOT CHANGE ANY SERVER DETAILS OR DB URLS IN THE SOURCE CODE.
// CHANGE THE CONFIG FILE ONLY.
package config

type AppConfig struct {
	TCPSocketConfig    *TCPSocketConfig   `yaml:"tcp"`
	Mongoconfig        *Mongoconfig       `yaml:"mongo"`
}

type Mongoconfig struct {
	URL         string                 `yaml:"url"`
	DBName 		string 				   `yaml:"db"`
	RawDB 		string 				   `yaml:"rawDB"`
	Collections *MongoCollectionConfig `yaml:"collections"`
}

type MongoCollectionConfig struct {
	LocationHistoriesCollection string `yaml:"location_histories"`
	VehicleDetailsCollection    string `yaml:"vehicle_details"`
	FenceDetailsCollection  	string `yaml:"fence_details"`
}

type TCPSocketConfig struct {
	ServerURL string `yaml:"serverURL"`
	Port      string `yaml:"port"`
}
