package main

import (
	"fmt"
	"github.com/axhixh/rabbitmq-experiments/stream"
	"github.com/streadway/amqp"
	"log"
)

func errorHandler(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	log.Printf("Receiving messages from: ")
	url, err := stream.GetRabbitMQ()
	errorHandler(err, "Unable to get URL for RabbitMQ")

	fmt.Printf("%s\n", url)
	conn, err := amqp.Dial(url)
	errorHandler(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	errorHandler(err, "Unable to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("q1", false, false, false, false, nil)
	errorHandler(err, "Unable to create queue")

	msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	errorHandler(err, "Failed to register consumer")

	forever := make(chan bool)

	go func() {
		for d := range msg {
			log.Printf("Msg: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting. Press CTRL+C to exit")
	<-forever
}
