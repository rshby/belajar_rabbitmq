package main

import (
	"belajar_rabbitmq/topic/publisher/rabbit"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	if err := godotenv.Load("./../../.env"); err != nil {
		logrus.Fatalf("cant load env file : %v", err)
	}
}

func main() {
	// connect rabbitmq
	conn := rabbit.ConnectRabbitMQ()
	defer conn.Close()

	// open channel
	ch, err := conn.Channel()
	defer ch.Close()
	if err != nil {
		logrus.Fatalf("cant open channel : %v", err)
	}

	// declare exchange
	exchangeName := "aggrements"
	if err := ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil); err != nil {
		logrus.Fatalf("cant declare exchange aggrements : %v", err)
	}

	// declare queue
	berlinQueue, err := ch.QueueDeclare("berlin_aggrement", true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("cant declare queue berlin_aggrement : %v", err)
	}

	allQueue, err := ch.QueueDeclare("all_aggrement", true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("cant declare queue all_aggrement : %v", err)
	}

	storeQueue, err := ch.QueueDeclare("store_aggrement", true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("cant declare queue store_aggrement : %v", err)
	}

	// declare binding
	berlinKey := "aggrements.eu.berlin.#"
	if err := ch.QueueBind(berlinQueue.Name, berlinKey, exchangeName, false, nil); err != nil {
		logrus.Fatalf("cant bind exchange to berlin_aggrement queue : %v", err)
	}

	allKey := "aggrements.#"
	if err := ch.QueueBind(allQueue.Name, allKey, exchangeName, false, nil); err != nil {
		logrus.Fatalf("cant bind exchange to all_aggrement queue : %v", err)
	}

	storeKey := "aggrements.eu.*.store"
	if err := ch.QueueBind(storeQueue.Name, storeKey, exchangeName, false, nil); err != nil {
		logrus.Fatalf("cant bind exchange to store_aggrement queue : %v", err)
	}

	// publish message
	keySend := "aggrements.eu.test.store"
	if err := ch.Publish(exchangeName, keySend, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(fmt.Sprintf("test %v", keySend)),
	}); err != nil {
		logrus.Fatalf("cant send message with key %v : %v", keySend, err)
	}

	logrus.Info(fmt.Sprintf("success send message to exchange %v, with routing key '%v'", exchangeName, keySend))
	fmt.Println("services sends all messages...")
}
