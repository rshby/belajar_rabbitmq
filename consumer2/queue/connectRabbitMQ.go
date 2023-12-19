package queue

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"log"
	"os"
)

// function to connect with RabbitMQ Server
func ConnectRabbitMQ() *amqp.Connection {
	dsn := fmt.Sprintf("amqp://%v:%v@localhost:%v",
		os.Getenv("RABBIT_USERNAME"),
		os.Getenv("RABBIT_PASSWORD"),
		os.Getenv("RABBIT_PORT"))

	conn, err := amqp.Dial(dsn)
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant dial RabbitMQ : %v", err.Error())
	}

	return conn
}
