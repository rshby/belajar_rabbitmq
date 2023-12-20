package main

import (
	"belajar_rabbitmq/direct/publisher/rabbit"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"sync"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	if err := godotenv.Load("./../../.env"); err != nil {
		logrus.Fatalf("error cant load env file : %v", err)
	}
}

func main() {
	// connect to rabbitMQ server
	conn := rabbit.ConnectRabbitMQ()
	defer conn.Close()

	// declare channel
	ch, err := conn.Channel()
	defer ch.Close()
	if err != nil {
		logrus.Fatalf("cant create channel : %v", err)
	}

	// declare exchange 'action_data' type direct
	exchangeName := "action_data"
	if err := ch.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil); err != nil {
		logrus.Fatalf("cant create exchange : %v", err)
	}

	// create queue 'create_data'
	createDataQueue, err := ch.QueueDeclare("create_data", true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("cant declare queue create_data : %v", err)
	}

	// create queue 'log_data'
	logDataQueue, err := ch.QueueDeclare("log_data", true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("cant declare queue log_data : %v", err)
	}

	// binding 'create_key' from exchange to queue 'create_data'
	if err := ch.QueueBind(createDataQueue.Name, "create_key", exchangeName, false, nil); err != nil {
		logrus.Fatalf("cant binding key 'create_key' : %v", err)
	}

	// binding 'log_key' from exchange to queue 'log_data'
	if err := ch.QueueBind(logDataQueue.Name, "log_key", exchangeName, false, nil); err != nil {
		logrus.Fatalf("cant binding key 'log_key' : %v", err)
	}

	wg := &sync.WaitGroup{}

	// send message to exchange with routing-key create_key
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		if err := ch.Publish(exchangeName, "create_key", false, false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte("test create data"),
			}); err != nil {
			logrus.Fatalf("cant send message to queue create_data : %v", err)
		}

		logrus.Info("send message to create_data")
	}(wg)

	// send message to exchange with routing-key log_key
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		if err := ch.Publish(exchangeName, "log_key", false, false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte("test log data"),
			}); err != nil {
			logrus.Fatalf("cant send message to queue log_data : %v", err)
		}

		logrus.Info("send message to log_data")
	}(wg)

	wg.Wait()
	fmt.Println("publisher sends all messages")
}
