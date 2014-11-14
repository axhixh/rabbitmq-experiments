package stream

import (
	"flag"
	"fmt"
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
