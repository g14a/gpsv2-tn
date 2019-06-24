package amqputils

import (
	"github.com/streadway/amqp"
	"gitlab.com/gpsv2/config"
	"gitlab.com/gpsv2/errorcheck"
	"sync"
)

var (
	amqpConn amqp.Connection
	amqpOnce sync.Once
)

func GetAMQPInstance() *amqp.Connection {
	amqpOnce.Do(func() {
		getAMQPConnection()
	})

	return &amqpConn
}

func getAMQPConnection() {
	amqpURL := config.GetAppConfig().AMQPConfig.URL

	conn, err := amqp.Dial(amqpURL)
	errorcheck.CheckError(err)

	amqpConn = *conn
}
