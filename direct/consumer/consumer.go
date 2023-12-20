package main

import (
	"belajar_rabbitmq/direct/consumer/rabbit"
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
	// connect to rabbit
	conn := rabbit.ConnectRabbitMQ()
	defer conn.Close()

	// open channel
	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("cant open channel %v", err)
	}

	// consume queue
	msgCreateData, err := ch.Consume("create_data", "", true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("cant consume create_data queue : %v", err)
	}

	msgLogData, err := ch.Consume("log_data", "", true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("cant consume log_data queue : %v", err)
	}

	forever := make(chan bool)

	// consume create_date
	go func() {
		for message := range msgCreateData {
			logrus.WithField("action", "create").Info(fmt.Sprintf("receive message create : %v", string(message.Body)))
		}
	}()

	// consume log_data
	go func() {
		for message := range msgLogData {
			logrus.WithField("action", "log").Info(fmt.Sprintf("receive message log : %v", string(message.Body)))
		}
	}()

	fmt.Println("service has been waiting messages...")
	<-forever
}
