package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	fmt.Println("Receiving message")
	conn, err := amqp.Dial("amqp://admin:VgkDl7PM5DzY@172.17.0.3:5672/")
	handleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	handleError(err, "Unable to open a channel")
	defer ch.Close()

    err = ch.ExchangeDeclare("chash-x", "x-consistent-hash", true, false, false, false, nil)
    handleError(err, "Unable to declare exchange")

	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	handleError(err, "Unable to create queue")

    err = ch.QueueBind(q.Name, "100", "chash-x", false, nil)
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
