package main

import (
	"belajar_rabbitmq/topic/consumer/rabbit"
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
	// connect to rabbit
	conn := rabbit.ConnectRabbit()
	defer conn.Close()

	// open a channel
	ch, err := conn.Channel()
	defer ch.Close()
	if err != nil {
		logrus.Fatalf("cant open a channel : %v", err)
	}

	// consume a queue
	berlinQueue := "berlin_aggrement"
	berlinMsg, err := ch.Consume(berlinQueue, "", true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("cant consume queue %v : %v\n", berlinQueue, err)
	}

	allAggrementQueue := "all_aggrement"
	allMsg, err := ch.Consume(allAggrementQueue, "", true, false, false, false, amqp.Table{
		"x-priority-max": "0-256",
	})
	if err != nil {
		logrus.Fatalf("cant consume queue %v : %v\n", allAggrementQueue, err)
	}

	storeQueue := "store_aggrement"
	storeMsg, err := ch.Consume(storeQueue, "", true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("cant consume queue :%v : %v\n", storeQueue, err)
	}

	forever := make(chan bool)

	// consume berlin_aggrement
	go func() {
		for message := range berlinMsg {
			logrus.WithField("queue", berlinQueue).Info(fmt.Sprintf("receive message : %v", string(message.Body)))
		}
	}()

	// consume all_aggrement
	go func() {
		for message := range allMsg {
			logrus.WithField("queue", allAggrementQueue).Info(fmt.Sprintf("receive message : %v", string(message.Body)))
		}
	}()

	// consume store_aggrement
	go func() {
		for message := range storeMsg {
			logrus.WithField("queue", storeQueue).Info(fmt.Sprintf("receive message : %v", string(message.Body)))
		}
	}()

	fmt.Println("service is waiting a message...")
	<-forever
}
