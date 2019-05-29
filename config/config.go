package config

type AppConfig struct {
	Mongoconfig     *Mongoconfig     `yaml:"mongo"`
	TCPSocketConfig *TCPSocketConfig `yaml:"tcp"`
}

type Mongoconfig struct {
	URL         string
	Collections *MongoCollectionConfig `yaml:"collections"`
}

type MongoCollectionConfig struct {
	LocationHistoriesCollection string `yaml:"location_histories"`
	VehicleDetailsCollection    string `yaml:"vehicle_details"`
}

type TCPSocketConfig struct {
	ServerURL string `yaml:"serverURL"`
	Port      string `yaml:"port"`
}
