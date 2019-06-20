// Package config is the run time representation of
// the config.yml file and takes in all the details
// of the databases and tcp server and port.
package config

type AppConfig struct {
	Mongoconfig        *Mongoconfig       `yaml:"mongo"`
	HistoryMongoConfig *BackupMongoConfig `yaml:"mongohistory"`
	TCPSocketConfig    *TCPSocketConfig   `yaml:"tcp"`
}

type Mongoconfig struct {
	URL         string                 `yaml:"url"`
	Collections *MongoCollectionConfig `yaml:"collections"`
}

type MongoCollectionConfig struct {
	LocationHistoriesCollection string `yaml:"location_histories"`
	VehicleDetailsCollection    string `yaml:"vehicle_details"`
}

type BackupMongoConfig struct {
	BackupURL         string                        `yaml:"backupurl"`
	BackupCollections *MongoHistoryCollectionConfig `yaml:"backupcollections"`
}

type MongoHistoryCollectionConfig struct {
	BackupLocationHistoriesColl string `yaml:"location_histories"`
	RawDataCollection           string `yaml:"raw_data"`
}

type TCPSocketConfig struct {
	ServerURL string `yaml:"serverURL"`
	Port      string `yaml:"port"`
}
