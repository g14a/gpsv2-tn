// Package config is the run time representation of
// the config.yml file and takes in all the details
// of the databases and tcp server and port
package config

type AppConfig struct {
	AMQPConfig         *AMQPConfig        `yaml:"amqp"`
	TCPSocketConfig    *TCPSocketConfig   `yaml:"tcp"`
}

type AMQPConfig struct {
	URL       string `yaml:"amqpurl"`
	AMQPQueue string `yaml:"queuename"`
}

type TCPSocketConfig struct {
	ServerURL string `yaml:"serverURL"`
	Port      string `yaml:"port"`
}
