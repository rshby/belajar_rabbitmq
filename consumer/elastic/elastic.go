package elastic

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func ConnectToElastic() *elasticsearch.Client {
	config := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://localhost:%v", os.Getenv("ELASTIC_PORT")),
		},
	}

	client, err := elasticsearch.NewClient(config)
	if err != nil {
		logrus.Error(err.Error())
		log.Fatalf("error cant connect to elastic : %v", err.Error())
	}

	return client
}
