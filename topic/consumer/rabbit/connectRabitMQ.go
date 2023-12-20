package rabbit

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)

// func connect to RabbitMQ
func ConnectRabbit() *amqp.Connection {
	dsn := fmt.Sprintf("amqp://%v:%v@%v:%v",
		os.Getenv("RABBIT_USERNAME"),
		os.Getenv("RABBIT_PASSWORD"),
		os.Getenv("RABBIT_HOST"),
		os.Getenv("RABBIT_PORT"))

	conn, err := amqp.Dial(dsn)
	if err != nil {
		logrus.Fatalf("cant connect to RabbitMQ : %v", err)
	}

	return conn
}
