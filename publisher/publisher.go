package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"log"
	"os"
	"sync"
	"time"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	if err := godotenv.Load("./../.env"); err != nil {
		log.Fatalf("cant load env file : %v", err.Error())
	}
}

func main() {
	dnsRabbit := fmt.Sprintf("amqp://%v:%v@localhost:%v/",
		os.Getenv("RABBIT_USERNAME"),
		os.Getenv("RABBIT_PASSWORD"),
		os.Getenv("RABBIT_PORT"),
	)

	conn, err := amqp.Dial(dnsRabbit)
	defer conn.Close()
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error connection rabbit : %v", err.Error())
	}

	ch, err := conn.Channel()
	defer ch.Close()
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error create channel : %v", err.Error())
	}

	q, err := ch.QueueDeclare("golang-queue", false, false, false, false, nil)
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error create queue declare : %v", err.Error())
	}

	queueLog, err := ch.QueueDeclare("golang-log", false, false, false, false, nil)
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error create queue declare : %v", err.Error())
	}

	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}
	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup, mtx *sync.Mutex, i int) {
			wg.Add(1)
			mtx.Lock()
			defer func() {
				mtx.Unlock()
				wg.Done()
			}()

			body := fmt.Sprintf("message ke-%v", i+1)
			_ = ch.Publish("", q.Name, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})

			logrus.Info(fmt.Sprintf("success send message ke-%v", i+1))
		}(wg, mtx, i)

		go func(wg *sync.WaitGroup, mtx *sync.Mutex, i int) {
			wg.Add(1)
			mtx.Lock()
			defer func() {
				mtx.Unlock()
				wg.Done()
			}()
			body := map[string]any{
				"insert_at": time.Now().UTC().Add(7 * time.Hour).Format("2006-01-02 15:04:05"),
				"request":   fmt.Sprintf("%v", i+31),
				"response":  fmt.Sprintf("get data id %v", i+31),
			}

			jsonMessage, _ := json.Marshal(&body)

			_ = ch.Publish("", queueLog.Name, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         jsonMessage,
			})

			logrus.Info(fmt.Sprintf("success send Log with id : %v - %v", i+11, body))
		}(wg, mtx, i)
	}

	wg.Wait()
	logrus.Info("success send all messages")
}
