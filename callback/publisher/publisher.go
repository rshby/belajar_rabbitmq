package main

import (
	"belajar_rabbitmq/callback/publisher/rabbit"
	"belajar_rabbitmq/consumer/helper"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"strconv"
	"sync"
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
	conn := rabbit.ConnectRabbitMQ()
	defer conn.Close()

	// open channel
	ch, err := conn.Channel()
	defer ch.Close()
	if err != nil {
		logrus.Fatalf("cant open channel : %v", err)
	}

	// declare queue
	fibonaciQueueName := "fibonaci"
	fibonaciQueue, err := ch.QueueDeclare(fibonaciQueueName, true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("cant declare queue '%v' : %v", fibonaciQueueName, err)
	}

	// declare random queue to store response
	callbackQueue, err := ch.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		logrus.Fatalf("cant decalre queue callback : %v", err)
	}

	// consume queue callback
	msgCallback, err := ch.Consume(callbackQueue.Name, "", true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("cant consume callback queue : %v", err)
	}

	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	// publish a message
	n := 6
	for i := 0; i < 20000; i++ {
		go func(wg *sync.WaitGroup, mtx *sync.Mutex, i int) {
			wg.Add(1)
			mtx.Lock()
			defer func() {
				wg.Done()
				mtx.Unlock()
			}()

			corrId := helper.GenerateRandomID()
			if err := ch.Publish("", fibonaciQueue.Name, false, false,
				amqp.Publishing{
					DeliveryMode:  amqp.Persistent,
					ContentType:   "text/plain",
					CorrelationId: corrId,
					ReplyTo:       callbackQueue.Name,
					Body:          []byte(strconv.Itoa(n)),
				}); err != nil {
				logrus.Fatalf("cant send message : %v", err)
			}

			logrus.Info(fmt.Sprintf("success send message %v, i = %v", n, i))

			// consume callback -> receive response
			for response := range msgCallback {
				if response.CorrelationId == corrId {
					logrus.Info(fmt.Sprintf("receive response : %v, i = %v", string(response.Body), i))
					break
				}
			}
		}(wg, mtx, i)
	}

	wg.Wait()
	fmt.Println("service stopped")
}
