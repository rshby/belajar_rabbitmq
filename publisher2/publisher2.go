package main

import (
	"belajar_rabbitmq/publisher2/queue"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	if err := godotenv.Load("./../.env"); err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant load env : %v", err.Error())
	}
}

func main() {
	// connect to rabbit server
	conn := queue.ConnectRabbitMQ()
	defer conn.Close()

	// create channel
	ch, err := conn.Channel()
	defer ch.Close()
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant create channel : %v", err.Error())
	}

	// create exchange
	if err := ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil); err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant create exchange logs : %v", err.Error())
	}

	// decalre queue
	queue1, err := ch.QueueDeclare("queue1", true, false, false,
		false, nil)
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant create queue 1 : %v", err.Error())
	}

	queue2, err := ch.QueueDeclare("queue2", true, false, false,
		false, nil)
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant create queue2 : %v", err.Error())
	}

	// binding
	routingKey := "this_queue"
	exchangeName := "logs"
	if err := ch.QueueBind(queue1.Name, routingKey, exchangeName, false, nil); err != nil {
		logrus.Error(err.Error())
		log.Fatalf("cant bind queue1 : %v", err.Error())
	}

	if err := ch.QueueBind(queue2.Name, routingKey, exchangeName, false, nil); err != nil {
		logrus.Error(err.Error())
		log.Fatalf("cant bind queue2 : %v", err.Error())
	}

	body := []byte("Eko Saputra")
	// publish to exchange
	if err := ch.Publish(exchangeName, routingKey, false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // kalo fanout ga ngaruh
			ContentType:  "text/plain",
			Body:         body,
		}); err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant publish message : %v", err.Error())
	}

	// success send message
	logrus.Info(fmt.Sprintf("success send message : %v", string(body)))
}
