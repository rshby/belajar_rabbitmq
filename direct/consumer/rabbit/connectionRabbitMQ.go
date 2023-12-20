package rabbit

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)

// function to connection with RabbitMQ
func ConnectRabbitMQ() *amqp.Connection {
	dsn := fmt.Sprintf("amqp://%v:%v@%v:%v",
		os.Getenv("RABBIT_USERNAME"),
		os.Getenv("RABBIT_PASSWORD"),
		os.Getenv("RABBIT_HOST"),
		os.Getenv("RABBIT_PORT"))

	conn, err := amqp.Dial(dsn)
	if err != nil {
		logrus.Fatalf("cant connect to RabbitMQ server : %v", err)
	}

	// success connect to rabbit
	return conn
}
