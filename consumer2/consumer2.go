package main

import (
	"belajar_rabbitmq/consumer2/queue"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	if err := godotenv.Load("./../.env"); err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant load env file : %v", err.Error())
	}
}

func main() {
	// connect to RabbitMQ
	conn := queue.ConnectRabbitMQ()
	defer conn.Close()

	// declare channel
	ch, err := conn.Channel()
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant create channel : %v", err.Error())
	}

	// declare exchange
	exchangeName := "logs"
	if err := ch.ExchangeDeclare(exchangeName, "fanout", true, false, false, false, nil); err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant declare exchange logs : %v", err.Error())
	}

	// declare queue
	queue1, err := ch.QueueDeclare("queue1", true, false, false, false, nil)
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("cant declare queue1 : %v", err.Error())
	}

	queue2, err := ch.QueueDeclare("queue2", true, false, false, false, nil)
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("cant declare queue2 : %v", err.Error())
	}

	// binding
	routingKey := "this_queue"
	if err := ch.QueueBind(queue1.Name, routingKey, exchangeName, false, nil); err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant binding %v : %v", queue1.Name, err.Error())
	}

	if err := ch.QueueBind(queue2.Name, routingKey, exchangeName, false, nil); err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant binding %v : %v", queue2.Name, err.Error())
	}

	// consume queue
	messageQueue1, err := ch.Consume(queue1.Name, "", true, false, false, false, nil)
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant consume %v : %v", queue1.Name, err.Error())
	}

	messageQueue2, err := ch.Consume(queue2.Name, "", true, false, false, false, nil)
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant consume %v : %v", queue2.Name, err.Error())
	}

	var forever = make(chan bool)
	go func() {
		for message := range messageQueue1 {
			response := fmt.Sprintf("receive message : %v", string(message.Body))
			logrus.WithField("queue", queue1.Name).Info(response)
		}
	}()

	go func() {
		for message := range messageQueue2 {
			response := fmt.Sprintf("receive message : %v", string(message.Body))
			logrus.WithField("queue", queue2.Name).Info(response)
		}
	}()

	fmt.Println("waiting message....")
	<-forever
}
