package main

import (
	"belajar_rabbitmq/fanout/consumer/queue"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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
	// connection to rabbitMQ
	conn := queue.ConnectRabbitMQ()
	defer conn.Close()

	// create channel
	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("error cant create channel : %v", err)
	}
	defer ch.Close()

	// declare exchange
	exchangeName := "fanout1"
	if err := ch.ExchangeDeclare(exchangeName, "fanout", true, false, false, false, nil); err != nil {
		logrus.Fatalf("error cant declare exchange : %v", err)
	}

	// queue declare
	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		logrus.Fatalf("error cant decalre queue : %v", err)
	}

	// binding queue - exchange
	routingKey := "fanout1_key"
	if err := ch.QueueBind(q.Name, routingKey, exchangeName, false, nil); err != nil {
		logrus.Fatalf("error cant binding queue to exchange : %v", err)
	}

	// consume message
	msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("error cant consume queue : %v", err)
	}

	forever := make(chan bool)
	go func() {
		for message := range msg {
			response := fmt.Sprintf("receive message : %v", string(message.Body))
			logrus.Info(response)
		}
	}()

	fmt.Println("waiting message...")
	<-forever
}
