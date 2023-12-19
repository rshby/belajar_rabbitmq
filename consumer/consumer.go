package main

import (
	"belajar_rabbitmq/consumer/elastic"
	"belajar_rabbitmq/consumer/helper"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	if err := godotenv.Load("./../.env"); err != nil {
		logrus.Error(err.Error())
		log.Fatalf("cant load env file : %v", err.Error())
	}
}

func main() {
	esClient := elastic.ConnectToElastic()
	dnsRabbit := fmt.Sprintf("amqp://%v:%v@localhost:%v/",
		os.Getenv("RABBIT_USERNAME"),
		os.Getenv("rabbit_password"),
		os.Getenv("RABBIT_PORT"))

	conn, err := amqp.Dial(dnsRabbit)
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant dial rabbit server : %v", err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error create channel : %v", err.Error())
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"golang-queue", //name
		false,          // durable
		false,          //delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {

	}

	logTask, err := ch.QueueDeclare(
		"golang-log",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("cant decaler queue golang-log : %v", err.Error())
	}

	msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error create channel consume : %v", err.Error())
	}

	msgLog, err := ch.Consume(logTask.Name, "", true, false, false, false, nil)
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("cant consule golang-log :%v", err.Error())
	}

	forever := make(chan bool)
	go func() {
		for d := range msg {
			message := fmt.Sprintf("receive message : %v", string(d.Body))
			logrus.Info(message)
		}
	}()

	// method log and store in elastic
	go func() {
		for d := range msgLog {
			response := map[string]any{}
			_ = json.Unmarshal(d.Body, &response)
			res, err := esClient.Index(
				os.Getenv("ELASTIC_INDEX"),
				strings.NewReader(string(d.Body)),
				esClient.Index.WithDocumentID(helper.GenerateRandomID()),
				esClient.Index.WithRefresh("true"))
			if err != nil {
				logrus.Error(err.Error())
				log.Fatalf("error send to elastic : %v", err.Error())
			}
			defer res.Body.Close()

			if res.IsError() {
				logrus.Error(err.Error())
				log.Fatalf("error res isError : %v", err.Error())
			}

			logrus.Info(response)
		}
	}()

	logrus.Info("waiting for message")
	<-forever
}
