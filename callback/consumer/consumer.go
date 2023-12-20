package main

import (
	"belajar_rabbitmq/callback/consumer/helper"
	"belajar_rabbitmq/callback/consumer/rabbit"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"strconv"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	if err := godotenv.Load("./../../.env"); err != nil {
		logrus.Fatalf("cant load env file : %v", err)
	}
}

func main() {
	// connect to rabbitMQ
	conn := rabbit.ConnectToRabbitMQ()
	defer conn.Close()

	// open channel
	ch, err := conn.Channel()
	defer ch.Close()
	if err != nil {
		logrus.Fatalf("cant open channel : %v", err.Error())
	}

	// consume queue
	msgNumber, err := ch.Consume("fibonaci", "", false, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("cant consume fibonaci queue : %v", err.Error())
	}

	forever := make(chan bool)
	go func() {
		for message := range msgNumber {
			n, err := strconv.Atoi(string(message.Body))
			if err != nil {
				logrus.Fatalf("cant convert body : %v", err.Error())
			}

			fibonaci := helper.Fibonaci(n)
			logrus.Info(fmt.Sprintf("receive message : %v, with result fibonaci : %v", string(message.Body), fibonaci))

			// send to callback
			if err := ch.Publish("", message.ReplyTo, false, false,
				amqp.Publishing{
					DeliveryMode:  amqp.Persistent,
					ContentType:   "text/plain",
					CorrelationId: message.CorrelationId,
					Body:          []byte(strconv.Itoa(fibonaci)),
				}); err != nil {
				logrus.Fatalf("cant send response back : %v", err.Error())
			}

			logrus.Info(fmt.Sprintf("success send message callback : %v", fibonaci))
			_ = message.Ack(false)
		}
	}()

	fmt.Println("waiting message..")
	<-forever
}
