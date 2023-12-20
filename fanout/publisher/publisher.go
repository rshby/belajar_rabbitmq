package main

import (
	"belajar_rabbitmq/fanout/publisher/queue"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	if err := godotenv.Load("./../../.env"); err != nil {
		logrus.Fatalf("cant load file env : %v", err)
	}
}

func main() {
	// connect to rabbitMQ
	conn := queue.ConnectRabbitMQ()
	defer conn.Close()

	// create channel
	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("cant create channel : %v", err)
	}
	defer ch.Close()

	// declare exchange
	if err := ch.ExchangeDeclare("fanout1", "fanout", true, false, false, false, nil); err != nil {
		logrus.Fatalf("cant declare exchange : %v", err)
	}

	// publish to exchange
	if err := ch.Publish("fanout1", "", false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte("test fanout1"),
		}); err != nil {
		logrus.Fatalf("error cant send message : %v", err)
	}

	logrus.Info("success send message")
}
