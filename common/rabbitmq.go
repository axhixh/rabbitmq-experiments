package common

import (
	"flag"
	"fmt"
	"log"
)

func GetRabbitMQ() (string, error) {
	var url string
	flag.StringVar(&url, "mq", "", "URL for RabbitMQ [amqp://user:pass@server:5671]")
	flag.Parse()
	if url == "" {
		return "", fmt.Errorf("AMQP URL is missing [amqp://user:pass@server:5671/]")
	}
	return url, nil
}

func HandleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
