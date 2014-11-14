package main

import (
	"flag"
	"fmt"
	"github.com/axhixh/rabbitmq-experiments/stream"
	"github.com/streadway/amqp"
	"log"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func getQueueName() string {
	if flag.NArg() == 0 {
		panic("Please provide a queue name")
	}
	return flag.Arg(0)
}

func main() {
	log.Printf("Receiving message ")
	url, err := stream.GetRabbitMQ()
	handleError(err, "Unable to get RabbitMQ")

	queueName := getQueueName()

	conn, err := amqp.Dial(url)
	handleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	handleError(err, "Unable to open a channel")
	defer ch.Close()

	const exchangeName = "chash"
	err = ch.ExchangeDeclare(exchangeName, "x-consistent-hash", false, false, false, false, nil)
	handleError(err, "Unable to declare exchange")

	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	handleError(err, "Unable to create queue")

	err = ch.QueueBind(q.Name, "100", exchangeName, false, nil)
	msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	handleError(err, "Failed to register consumer")

	forever := make(chan bool)

	go func() {
		for d := range msg {
			log.Printf("%s: %s", q.Name, d.Body)
		}
	}()

	log.Printf("[*] Waiting. Press CTRL+C to exit")
	<-forever
}
